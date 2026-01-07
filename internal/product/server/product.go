package server

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/api/proto"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/internal/product/usecase"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	app *usecase.ProductService
	proto.ProductServiceServer
}

func NewProductServer(app *usecase.ProductService) *ProductServer {
	return &ProductServer{
		app: app,
	}
}

func (ps *ProductServer) Bind(server *grpc.Server) {
	proto.RegisterProductServiceServer(server, ps)
}

func (ps *ProductServer) BatchRequest(ctx context.Context, in *proto.BatchProductRequest) (*proto.BatchResponse, error) {
	var response proto.BatchResponse
	requests := in.GetRequests()
	batch := make([]*usecase.ProductRequest, len(requests))
	for i, req := range requests {
		r, err := toProductRequest(req)
		if err != nil {
			return nil, err
		}
		batch[i] = r
	}
	appResponse, err := ps.app.BatchRequest(ctx, batch)
	if err != nil {
		var daltyErr *daltyerrors.DaltyError
		if errors.As(err, &daltyErr) {
			return nil, protoerr.Handle(daltyErr)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	batchResponse := make([]*proto.Product, len(appResponse))
	for i, product := range appResponse {
		batchResponse[i] = toProtoProduct(product)
	}
	response.SetProducts(batchResponse)
	return &response, nil
}

func toProductRequest(in *proto.ProductRequest) (*usecase.ProductRequest, error) {
	if !in.HasFnrec() && !in.HasIntegrationId() {
		return nil, protoerr.InvalidArgument("request does not have any identifier",
			&protoerr.ValidationError{
				Message: "one of the following fields must be set",
				Members: []string{"fnrec", "integration_id"},
			})
	}
	var req usecase.ProductRequest
	if in.HasIntegrationId() {
		req.IntegrationID = in.GetIntegrationId()
		return &req, nil
	}
	req.Fnrec = in.GetFnrec()
	return &req, nil
}

func toProtoProduct(in *core.Product) *proto.Product {
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

func toProtoProductionType(in core.ProductionType) proto.ProductionType {
	switch in {
	case core.ProductionTypeProducing:
		return proto.ProductionType_PRODUCTION_TYPE_PRODUCING
	case core.ProductionTypePurchasing:
		return proto.ProductionType_PRODUCTION_TYPE_PURCHASING
	default:
		return proto.ProductionType_PRODUCTION_TYPE_UNKNOWN
	}
}
func toProtoProductGroup(pg core.ProductGroup) proto.ProductGroup {
	switch pg {
	case core.ProductGroupKitchens:
		return proto.ProductGroup_PRODUCT_GROUP_KITCHENS
	case core.ProductGroupCaseFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE
	case core.ProductGroupBeddingSets:
		return proto.ProductGroup_PRODUCT_GROUP_BEDDING_SETS
	case core.ProductGroupSofas:
		return proto.ProductGroup_PRODUCT_GROUP_SOFAS
	case core.ProductGroupCovers:
		return proto.ProductGroup_PRODUCT_GROUP_COVERS
	case core.ProductGroupBlankets:
		return proto.ProductGroup_PRODUCT_GROUP_BLANKETS
	case core.ProductGroupBedBasesWithStorage:
		return proto.ProductGroup_PRODUCT_GROUP_BED_BASES_WITH_STORAGE
	case core.ProductGroupSofaComponents:
		return proto.ProductGroup_PRODUCT_GROUP_SOFA_COMPONENTS
	case core.ProductGroupErgomotion:
		return proto.ProductGroup_PRODUCT_GROUP_ERGOMOTION
	case core.ProductGroupNonProducts:
		return proto.ProductGroup_PRODUCT_GROUP_NON_PRODUCTS
	case core.ProductGroupSmallFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE
	case core.ProductGroupMattresses:
		return proto.ProductGroup_PRODUCT_GROUP_MATTRESSES
	case core.ProductGroupSlattedBases:
		return proto.ProductGroup_PRODUCT_GROUP_SLATTED_BASES
	case core.ProductGroupMattressToppers:
		return proto.ProductGroup_PRODUCT_GROUP_MATTRESS_TOPPERS
	case core.ProductGroupPillows:
		return proto.ProductGroup_PRODUCT_GROUP_PILLOWS
	case core.ProductGroupBeds:
		return proto.ProductGroup_PRODUCT_GROUP_BEDS
	case core.ProductGroupBedBases:
		return proto.ProductGroup_PRODUCT_GROUP_BED_BASES
	case core.ProductGroupMiscellaneous:
		return proto.ProductGroup_PRODUCT_GROUP_MISCELLANEOUS
	case core.ProductGroupWardrobes:
		return proto.ProductGroup_PRODUCT_GROUP_WARDROBES
	case core.ProductGroupTextiles:
		return proto.ProductGroup_PRODUCT_GROUP_TEXTILES
	case core.ProductGroupElectronics:
		return proto.ProductGroup_PRODUCT_GROUP_ELECTRONICS
	case core.ProductGroupClothing:
		return proto.ProductGroup_PRODUCT_GROUP_CLOTHING
	case core.ProductGroupOrthopedics:
		return proto.ProductGroup_PRODUCT_GROUP_ORTHOPEDICS
	case core.ProductGroupCoffeeTables:
		return proto.ProductGroup_PRODUCT_GROUP_COFFEE_TABLES
	case core.ProductGroupHomeOffice:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_OFFICE
	case core.ProductGroupLivingRooms:
		return proto.ProductGroup_PRODUCT_GROUP_LIVING_ROOMS
	case core.ProductGroupLighting:
		return proto.ProductGroup_PRODUCT_GROUP_LIGHTING
	case core.ProductGroupDecor:
		return proto.ProductGroup_PRODUCT_GROUP_DECOR
	case core.ProductGroupSpaceOrganizationUpper:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case core.ProductGroupHomeCareUpper:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE_UPPER
	case core.ProductGroupSpaceOrganization:
		return proto.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION
	case core.ProductGroupHomeCare:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case core.ProductGroupHallways:
		return proto.ProductGroup_PRODUCT_GROUP_HALLWAYS
	case core.ProductGroupFurnitureProtectionAndCare:
		return proto.ProductGroup_PRODUCT_GROUP_FURNITURE_PROTECTION_AND_CARE
	case core.ProductGroupOutdoorFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_OUTDOOR_FURNITURE
	case core.ProductGroupStorage:
		return proto.ProductGroup_PRODUCT_GROUP_STORAGE
	case core.ProductGroupInterior:
		return proto.ProductGroup_PRODUCT_GROUP_INTERIOR
	case core.ProductGroupSeasonalProducts:
		return proto.ProductGroup_PRODUCT_GROUP_SEASONAL_PRODUCTS
	case core.ProductGroupFragrances:
		return proto.ProductGroup_PRODUCT_GROUP_FRAGRANCES

	// accessories & additional groups
	case core.ProductGroupMurphyBeds:
		return proto.ProductGroup_PRODUCT_GROUP_MURPHY_BEDS
	case core.ProductGroupBedAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_BED_ACCESSORIES
	case core.ProductGroupMurphyBedAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_MURPHY_BED_ACCESSORIES
	case core.ProductGroupSmallFurnitureAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE_ACCESSORIES
	case core.ProductGroupWardrobeAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_WARDROBE_ACCESSORIES
	case core.ProductGroupInteriorDecoration:
		return proto.ProductGroup_PRODUCT_GROUP_INTERIOR_DECORATION
	case core.ProductGroupSleepTherapy:
		return proto.ProductGroup_PRODUCT_GROUP_SLEEP_THERAPY
	case core.ProductGroupKingKoil:
		return proto.ProductGroup_PRODUCT_GROUP_KING_KOIL
	case core.ProductGroupErgomotionAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_ERGOMOTION_ACCESSORIES
	case core.ProductGroupChildrenBedBases:
		return proto.ProductGroup_PRODUCT_GROUP_CHILDREN_BED_BASES
	case core.ProductGroupPillowCovers:
		return proto.ProductGroup_PRODUCT_GROUP_PILLOW_COVERS
	case core.ProductGroupTableware:
		return proto.ProductGroup_PRODUCT_GROUP_TABLEWARE
	case core.ProductGroupSets:
		return proto.ProductGroup_PRODUCT_GROUP_SETS
	case core.ProductGroupChildrenBedrooms:
		return proto.ProductGroup_PRODUCT_GROUP_CHILDREN_BEDROOMS
	case core.ProductGroupSpaceOrganizationStorage:
		return proto.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION_STORAGE
	case core.ProductGroupBathroomProducts:
		return proto.ProductGroup_PRODUCT_GROUP_BATHROOM_PRODUCTS
	case core.ProductGroupToys:
		return proto.ProductGroup_PRODUCT_GROUP_TOYS
	case core.ProductGroupAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_ACCESSORIES
	case core.ProductGroupNewYear:
		return proto.ProductGroup_PRODUCT_GROUP_NEW_YEAR
	case core.ProductGroupArmchairs:
		return proto.ProductGroup_PRODUCT_GROUP_ARMCHAIRS
	case core.ProductGroupMassageChairs:
		return proto.ProductGroup_PRODUCT_GROUP_MASSAGE_CHAIRS
	case core.ProductGroupKafkaTest:
		return proto.ProductGroup_PRODUCT_GROUP_KAFKA_TEST
	case core.ProductGroupCaseFurnitureAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE_ACCESSORIES
	default:
		return proto.ProductGroup_PRODUCT_GROUP_UNSPECIFIED
	}
}

func toError(err error) error {
	if errors.Is(err, daltyerrors.ErrNotFound) {

	}
	return protoerr.InternalError(err)
}
