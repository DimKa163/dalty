package integration_test

import (
	"context"
	"fmt"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"

	"github.com/DimKa163/dalty/api/proto"
	"github.com/DimKa163/dalty/app/product"
	"github.com/DimKa163/dalty/internal/product/persistence"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDirectSpecification(t *testing.T) {
	ctx := context.Background()
	productIntegrationId := guid.NewString()
	container, server, err := run(ctx, t, func(ctx context.Context, dbPool *pgxpool.Pool) error {
		id, err := insertProduct(ctx, dbPool, "bed with bases", productIntegrationId, daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		leftId, err := insertProduct(ctx, dbPool, "bed", guid.NewString(), daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		rightId, err := insertProduct(ctx, dbPool, "base", guid.NewString(), daltymodel.ProductGroupBedBases, false, false)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, leftId, 1)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, rightId, 2)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to run test container or configure service: %v", err)
	}
	defer container.Terminate(ctx)
	defer server.Shutdown(ctx)
	// TODO add client and test cases
	conn, err := grpc.NewClient(server.Config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail configure client")
	}
	specClient := proto.NewSpecificationServiceClient(conn)
	var req proto.SpecificationRequest
	quantity := 2
	var specLineReq proto.SpecificationLine
	specLineReq.SetId(guid.NewString())
	specLineReq.SetIntegration(productIntegrationId)
	specLineReq.SetQuantity(int32(quantity))
	req.SetSpecificationLines([]*proto.SpecificationLine{&specLineReq})
	resp, err := specClient.Execute(ctx, &req)

	assert.NoError(t, err, "")
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetSpecifications())
	for _, spec := range resp.GetSpecifications() {
		assert.Equal(t, proto.SpecificationType_DIRECT_SPECIFICATION, spec.GetType())
		assert.Equal(t, proto.PickupStrategy_NEAREST, spec.GetStrategy())
		assert.NotEmpty(t, spec.GetChildProduct())
		for _, s := range spec.GetChildProduct() {
			switch s.GetProduct().GetGroup() {
			case proto.ProductGroup_PRODUCT_GROUP_BEDS:
				assert.Equal(t, proto.PickupStrategy_NEAREST, s.GetStrategy())
			case proto.ProductGroup_PRODUCT_GROUP_BED_BASES:
				assert.Equal(t, proto.PickupStrategy_FARTHEST, s.GetStrategy(), "pickup strategy should be "+
					"farthest for bed bases in compound specification")
				assert.Equal(t, int32(4), s.GetQuantity())
			}
		}
	}
}

func TestReverseSpecification(t *testing.T) {
	ctx := context.Background()
	leftProductIntegrationId := guid.NewString()
	rightProductIntegrationId := guid.NewString()
	container, server, err := run(ctx, t, func(ctx context.Context, dbPool *pgxpool.Pool) error {
		id, err := insertProduct(ctx, dbPool, "bed with bases", guid.NewString(), daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		leftId, err := insertProduct(ctx, dbPool, "bed", leftProductIntegrationId, daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		rightId, err := insertProduct(ctx, dbPool, "base", rightProductIntegrationId, daltymodel.ProductGroupBedBases, false, false)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, leftId, 1)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, rightId, 2)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to run test container or configure service: %v", err)
	}
	defer container.Terminate(ctx)
	defer server.Shutdown(ctx)

	conn, err := grpc.NewClient(server.Config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail configure client")
	}
	specClient := proto.NewSpecificationServiceClient(conn)

	var req proto.SpecificationRequest
	quantity := 2

	line1 := guid.NewString()
	line2 := guid.NewString()

	var specLineReq1 proto.SpecificationLine
	specLineReq1.SetId(line1)
	specLineReq1.SetIntegration(leftProductIntegrationId)
	specLineReq1.SetQuantity(int32(quantity))
	specLineReq1.SetRelateToId(line2)

	var specLineReq2 proto.SpecificationLine
	specLineReq2.SetId(line2)
	specLineReq2.SetIntegration(rightProductIntegrationId)
	specLineReq2.SetQuantity(int32(quantity))
	specLineReq2.SetRelateToId(line1)
	req.SetSpecificationLines([]*proto.SpecificationLine{&specLineReq1, &specLineReq2})

	resp, err := specClient.Execute(ctx, &req)

	assert.NoError(t, err, "")
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetSpecifications())
	for _, spec := range resp.GetSpecifications() {
		assert.Equal(t, proto.SpecificationType_REVERSE_SPECIFICATION, spec.GetType())
		assert.Equal(t, proto.PickupStrategy_NEAREST, spec.GetStrategy())
		assert.NotEmpty(t, spec.GetChildProduct())
		for _, s := range spec.GetChildProduct() {
			assert.Equal(t, int32(2), s.GetQuantity())
			switch s.GetProduct().GetGroup() {
			case proto.ProductGroup_PRODUCT_GROUP_BEDS:
				assert.Equal(t, proto.PickupStrategy_NEAREST, s.GetStrategy())
			case proto.ProductGroup_PRODUCT_GROUP_BED_BASES:
				assert.Equal(t, proto.PickupStrategy_FARTHEST, s.GetStrategy())
			}
		}
	}
}

func TestDefaultSpecification(t *testing.T) {
	ctx := context.Background()
	leftProductIntegrationId := guid.NewString()
	rightProductIntegrationId := guid.NewString()
	container, server, err := run(ctx, t, func(ctx context.Context, dbPool *pgxpool.Pool) error {
		id, err := insertProduct(ctx, dbPool, "bed with bases", guid.NewString(), daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		leftId, err := insertProduct(ctx, dbPool, "bed", leftProductIntegrationId, daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		rightId, err := insertProduct(ctx, dbPool, "base", rightProductIntegrationId, daltymodel.ProductGroupBedBases, false, false)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, leftId, 1)
		if err != nil {
			return err
		}
		_, err = insertRelation(ctx, dbPool, id, rightId, 2)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to run test container or configure service: %v", err)
	}
	defer container.Terminate(ctx)
	defer server.Shutdown(ctx)

	conn, err := grpc.NewClient(server.Config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail configure client")
	}
	specClient := proto.NewSpecificationServiceClient(conn)

	var req proto.SpecificationRequest
	quantity := 2

	line1 := guid.NewString()
	line2 := guid.NewString()

	var specLineReq1 proto.SpecificationLine
	specLineReq1.SetId(line1)
	specLineReq1.SetIntegration(leftProductIntegrationId)
	specLineReq1.SetQuantity(int32(quantity))

	var specLineReq2 proto.SpecificationLine
	specLineReq2.SetId(line2)
	specLineReq2.SetIntegration(rightProductIntegrationId)
	specLineReq2.SetQuantity(int32(quantity))
	req.SetSpecificationLines([]*proto.SpecificationLine{&specLineReq1, &specLineReq2})

	resp, err := specClient.Execute(ctx, &req)

	assert.NoError(t, err, "")
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetSpecifications())
	for _, spec := range resp.GetSpecifications() {
		assert.Equal(t, proto.SpecificationType_DEFAULT, spec.GetType())
		assert.Equal(t, proto.PickupStrategy_NEAREST, spec.GetStrategy())
		assert.NotEmpty(t, spec.GetChildProduct())
		for _, s := range spec.GetChildProduct() {
			assert.Equal(t, int32(2), s.GetQuantity())
			switch s.GetProduct().GetGroup() {
			case proto.ProductGroup_PRODUCT_GROUP_BEDS:
				assert.Equal(t, proto.PickupStrategy_NEAREST, s.GetStrategy())
			case proto.ProductGroup_PRODUCT_GROUP_BED_BASES:
				assert.Equal(t, proto.PickupStrategy_NEAREST, s.GetStrategy())
			}
		}
	}
}

func TestArchiveSpecification(t *testing.T) {
	ctx := context.Background()
	integrationId := guid.NewString()
	container, server, err := run(ctx, t, func(ctx context.Context, dbPool *pgxpool.Pool) error {
		_, err := insertProduct(ctx, dbPool, "bed", integrationId, daltymodel.ProductGroupBeds, false, true)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to run test container or configure service: %v", err)
	}
	defer container.Terminate(ctx)
	defer server.Shutdown(ctx)
	conn, err := grpc.NewClient(server.Config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail configure client")
	}
	specClient := proto.NewSpecificationServiceClient(conn)

	var req proto.SpecificationRequest
	quantity := 2

	line1 := guid.NewString()

	var specLineReq proto.SpecificationLine
	specLineReq.SetId(line1)
	specLineReq.SetIntegration(integrationId)
	specLineReq.SetQuantity(int32(quantity))

	req.SetSpecificationLines([]*proto.SpecificationLine{&specLineReq})

	resp, err := specClient.Execute(ctx, &req)

	assert.Nil(t, resp)
	assert.NotNil(t, err, "should be error")

	st, ok := status.FromError(err)
	assert.True(t, ok, "should be ok")
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	for _, detail := range protoerr.Flatten(st) {
		assert.Equal(t, int32(daltyerrors.DaltyErrorCodeOutOfStock), detail.GetCode())
	}
}

func TestNotFoundSpecification(t *testing.T) {
	ctx := context.Background()
	integrationId := guid.NewString()
	container, server, err := run(ctx, t, func(ctx context.Context, dbPool *pgxpool.Pool) error {
		_, err := insertProduct(ctx, dbPool, "bed", guid.NewString(), daltymodel.ProductGroupBeds, false, false)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to run test container or configure service: %v", err)
	}
	defer container.Terminate(ctx)
	defer server.Shutdown(ctx)
	conn, err := grpc.NewClient(server.Config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail configure client")
	}
	specClient := proto.NewSpecificationServiceClient(conn)

	var req proto.SpecificationRequest
	quantity := 2

	line1 := guid.NewString()

	var specLineReq proto.SpecificationLine
	specLineReq.SetId(line1)
	specLineReq.SetIntegration(integrationId)
	specLineReq.SetQuantity(int32(quantity))

	req.SetSpecificationLines([]*proto.SpecificationLine{&specLineReq})

	resp, err := specClient.Execute(ctx, &req)

	assert.Nil(t, resp)
	assert.NotNil(t, err, "should be error")

	st, ok := status.FromError(err)
	assert.True(t, ok, "should be ok")
	assert.Equal(t, codes.NotFound, st.Code())
	for _, detail := range protoerr.Flatten(st) {
		assert.Equal(t, int32(daltyerrors.DaltyErrorCodeResourceNotFound), detail.GetCode())
	}
}

func run(ctx context.Context, t *testing.T, beforeFn func(ctx context.Context, dbPool *pgxpool.Pool) error) (testcontainers.Container, *product.Server, error) {
	dbName := "dalty_product"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").
			WithStartupTimeout(90 * time.Second),
	}
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, _ := dbContainer.Host(ctx)
	port, _ := dbContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/%s?sslmode=disable", host, port.Port(), dbName)
	t.Logf("postgres started at: %s", dsn)
	server := product.NewServer(&product.Config{
		Addr:     ":8080",
		Database: dsn,
	})
	if err := server.AddServices(); err != nil {
		return nil, nil, err
	}
	if err := server.AddLogging(); err != nil {
		return nil, nil, err
	}
	server.Map()
	go func() {
		if err := server.Run(); err != nil {
			t.Logf("server error: %v", err)
		}
	}()
	if err = persistence.Migrate(server.PgPool, "../migrations"); err != nil {
		return nil, nil, err
	}
	if err := beforeFn(ctx, server.ServiceContainer.PgPool); err != nil {
		return nil, nil, err
	}
	return dbContainer, server, nil
}

func insertProduct(ctx context.Context, pool *pgxpool.Pool, name, integrationId string, group daltymodel.ProductGroup, isService, isArchive bool) (string, error) {
	id := guid.NewString()
	fnrec := fmt.Sprintf("fnrec_%s", id)
	if name == "" {
		name = fmt.Sprintf("product_%s", id)
	}
	_, err := pool.Exec(ctx, `INSERT INTO public.product(
	id, 
	smr_fnrec, 
	name, 
	smr_is_service, 
	is_archive,  
	type_id, 
	nrb_type_production_id, 
	nrb_account_product_id, 
	ask_non_standart_category_id, 
	nrb_integration_id, 
	smr_product_group_flag_id)
	VALUES ($1, 
	$2, 
	$3, 
	$4, 
	$5, 
	$6,
	$7,
	$8,
	$9,
	$10,
	$11);`,
		id,
		fnrec,
		name,
		isService,
		isArchive,
		daltymodel.ProductTypeMaterialAsset,
		daltymodel.ProductionTypeProducing,
		guid.NewString(),
		guid.NewString(),
		integrationId,
		group)
	return id, err
}

func insertRelation(ctx context.Context, pool *pgxpool.Pool, leftId, rightId string, amount int32) (string, error) {
	id := guid.NewString()
	_, err := pool.Exec(ctx, `INSERT INTO public.nrb_related_product(
	id,
	nrb_product_sku_id,
	nrb_product_mv_id,
	nrb_amount_mv)
	VALUES ($1, $2, $3, $4);`, id, leftId, rightId, amount)
	return id, err
}
