package server

import (
	"context"
	"fmt"
	"github.com/DimKa163/dalty/api/proto"
	"github.com/DimKa163/dalty/internal/product/usecase"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"google.golang.org/grpc"
)

type SpecificationServer struct {
	app *usecase.SpecificationService
	proto.UnimplementedSpecificationServiceServer
}

func NewSpecificationServer(app *usecase.SpecificationService) *SpecificationServer {
	return &SpecificationServer{app: app}
}

func (ss *SpecificationServer) Bind(server *grpc.Server) {
	proto.RegisterSpecificationServiceServer(server, ss)
}

func (ss *SpecificationServer) Execute(ctx context.Context, in *proto.SpecificationRequest) (*proto.SpecificationResponse, error) {
	var response proto.SpecificationResponse
	var specReq usecase.SpecRequest
	var err error
	var res []*daltymodel.Specification
	if err = toIn(in, &specReq); err != nil {
		return nil, err
	}
	res, err = ss.app.Execute(ctx, &specReq)
	if err != nil {
		return nil, err
	}
	specs := make([]*proto.Specification, len(res))
	for i, r := range res {
		specs[i] = toSpecification(r)
	}
	response.SetSpecifications(specs)
	return &response, nil
}

func toIn(in *proto.SpecificationRequest, request *usecase.SpecRequest) (err error) {
	lines := in.GetSpecificationLines()
	if len(lines) == 0 {
		return protoerr.InvalidArgument("no specification lines", &protoerr.ValidationError{
			Message: "no specification lines",
			Members: []string{
				"specification_lines",
			},
		})
	}
	request.Specs = make([]*usecase.Spec, len(lines))

	for i, line := range lines {
		if !line.HasQuantity() {
			err = fmt.Errorf("invalid specification line: %w", err)
		}
		if !line.HasFnrec() && !line.HasIntegration() {
			err = fmt.Errorf("invalid specification line: %w", err)
		}
		request.Specs[i] = &usecase.Spec{
			IntegrationID: line.GetIntegration(),
			Fnrec:         line.GetFnrec(),
			Quantity:      line.GetQuantity(),
		}
	}
	return
}

func toSpecification(spec *daltymodel.Specification) *proto.Specification {
	var specification proto.Specification
	specification.SetProduct(toLine(spec.Product))
	specification.SetType(proto.SpecificationType(spec.Type))
	specification.SetStrategy(proto.PickupStrategy(spec.Strategy))
	childProducts := make([]*proto.Line, len(spec.ChildProducts))
	for i, childProduct := range spec.ChildProducts {
		childProducts[i] = toLine(childProduct)
	}
	specification.SetChildProduct(childProducts)
	return &specification
}

func toLine(ln *daltymodel.Line) *proto.Line {
	var line proto.Line
	line.SetProduct(toProtoProductV2(ln.Product))
	line.SetQuantity(ln.Quantity)
	line.SetStrategy(proto.PickupStrategy(ln.Strategy))
	return &line
}

func toProtoProductV2(in *daltymodel.Product) *proto.Product {
	var out proto.Product
	out.SetId(in.ID.String())
	out.SetName(in.Name)
	out.SetProductionType(toProtoProductionType(in.ProductionType))
	out.SetFnrec(in.Fnrec)
	out.SetIsService(in.IsService)
	out.SetGroup(toProtoProductGroup(in.Group))
	out.SetSeriesId(in.SeriesID)
	out.SetCategoryId(in.CategoryID)
	out.SetAccountProvider(in.AccountProviderId)
	out.SetNonStandardCategoryId(in.NonStandardCategory)
	var pack proto.Pack
	pack.SetLength(in.Length)
	pack.SetHeight(in.Height)
	pack.SetWidth(in.Width)
	pack.SetVolume(in.Volume)
	pack.SetWeight(in.Weight)
	out.SetPack(&pack)
	return &out
}
