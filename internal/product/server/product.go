package server

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/api/proto"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/internal/product/usecase"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"github.com/DimKa163/dalty/pkg/daltymodel"
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

func toProtoProductionType(in daltymodel.ProductionType) proto.ProductionType {
	switch in {
	case core.ProductionTypeProducing:
		return proto.ProductionType_PRODUCTION_TYPE_PRODUCING
	case core.ProductionTypePurchasing:
		return proto.ProductionType_PRODUCTION_TYPE_PURCHASING
	default:
		return proto.ProductionType_PRODUCTION_TYPE_UNKNOWN
	}
}
func toProtoProductGroup(pg daltymodel.ProductGroup) proto.ProductGroup {
	switch pg {
	case daltymodel.ProductGroupKitchens:
		return proto.ProductGroup_PRODUCT_GROUP_KITCHENS
	case daltymodel.ProductGroupCaseFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE
	case daltymodel.ProductGroupBeddingSets:
		return proto.ProductGroup_PRODUCT_GROUP_BEDDING_SETS
	case daltymodel.ProductGroupSofas:
		return proto.ProductGroup_PRODUCT_GROUP_SOFAS
	case daltymodel.ProductGroupCovers:
		return proto.ProductGroup_PRODUCT_GROUP_COVERS
	case daltymodel.ProductGroupBlankets:
		return proto.ProductGroup_PRODUCT_GROUP_BLANKETS
	case daltymodel.ProductGroupBedBasesWithStorage:
		return proto.ProductGroup_PRODUCT_GROUP_BED_BASES_WITH_STORAGE
	case daltymodel.ProductGroupSofaComponents:
		return proto.ProductGroup_PRODUCT_GROUP_SOFA_COMPONENTS
	case daltymodel.ProductGroupErgomotion:
		return proto.ProductGroup_PRODUCT_GROUP_ERGOMOTION
	case daltymodel.ProductGroupNonProducts:
		return proto.ProductGroup_PRODUCT_GROUP_NON_PRODUCTS
	case daltymodel.ProductGroupSmallFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE
	case daltymodel.ProductGroupMattresses:
		return proto.ProductGroup_PRODUCT_GROUP_MATTRESSES
	case daltymodel.ProductGroupSlattedBases:
		return proto.ProductGroup_PRODUCT_GROUP_SLATTED_BASES
	case daltymodel.ProductGroupMattressToppers:
		return proto.ProductGroup_PRODUCT_GROUP_MATTRESS_TOPPERS
	case daltymodel.ProductGroupPillows:
		return proto.ProductGroup_PRODUCT_GROUP_PILLOWS
	case daltymodel.ProductGroupBeds:
		return proto.ProductGroup_PRODUCT_GROUP_BEDS
	case daltymodel.ProductGroupBedBases:
		return proto.ProductGroup_PRODUCT_GROUP_BED_BASES
	case daltymodel.ProductGroupMiscellaneous:
		return proto.ProductGroup_PRODUCT_GROUP_MISCELLANEOUS
	case daltymodel.ProductGroupWardrobes:
		return proto.ProductGroup_PRODUCT_GROUP_WARDROBES
	case daltymodel.ProductGroupTextiles:
		return proto.ProductGroup_PRODUCT_GROUP_TEXTILES
	case daltymodel.ProductGroupElectronics:
		return proto.ProductGroup_PRODUCT_GROUP_ELECTRONICS
	case daltymodel.ProductGroupClothing:
		return proto.ProductGroup_PRODUCT_GROUP_CLOTHING
	case daltymodel.ProductGroupOrthopedics:
		return proto.ProductGroup_PRODUCT_GROUP_ORTHOPEDICS
	case daltymodel.ProductGroupCoffeeTables:
		return proto.ProductGroup_PRODUCT_GROUP_COFFEE_TABLES
	case daltymodel.ProductGroupHomeOffice:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_OFFICE
	case daltymodel.ProductGroupLivingRooms:
		return proto.ProductGroup_PRODUCT_GROUP_LIVING_ROOMS
	case daltymodel.ProductGroupLighting:
		return proto.ProductGroup_PRODUCT_GROUP_LIGHTING
	case daltymodel.ProductGroupDecor:
		return proto.ProductGroup_PRODUCT_GROUP_DECOR
	case daltymodel.ProductGroupSpaceOrganizationUpper:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case daltymodel.ProductGroupHomeCareUpper:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE_UPPER
	case daltymodel.ProductGroupSpaceOrganization:
		return proto.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION
	case daltymodel.ProductGroupHomeCare:
		return proto.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case daltymodel.ProductGroupHallways:
		return proto.ProductGroup_PRODUCT_GROUP_HALLWAYS
	case daltymodel.ProductGroupFurnitureProtectionAndCare:
		return proto.ProductGroup_PRODUCT_GROUP_FURNITURE_PROTECTION_AND_CARE
	case daltymodel.ProductGroupOutdoorFurniture:
		return proto.ProductGroup_PRODUCT_GROUP_OUTDOOR_FURNITURE
	case daltymodel.ProductGroupStorage:
		return proto.ProductGroup_PRODUCT_GROUP_STORAGE
	case daltymodel.ProductGroupInterior:
		return proto.ProductGroup_PRODUCT_GROUP_INTERIOR
	case daltymodel.ProductGroupSeasonalProducts:
		return proto.ProductGroup_PRODUCT_GROUP_SEASONAL_PRODUCTS
	case daltymodel.ProductGroupFragrances:
		return proto.ProductGroup_PRODUCT_GROUP_FRAGRANCES

	// accessories & additional groups
	case daltymodel.ProductGroupMurphyBeds:
		return proto.ProductGroup_PRODUCT_GROUP_MURPHY_BEDS
	case daltymodel.ProductGroupBedAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_BED_ACCESSORIES
	case daltymodel.ProductGroupMurphyBedAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_MURPHY_BED_ACCESSORIES
	case daltymodel.ProductGroupSmallFurnitureAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE_ACCESSORIES
	case daltymodel.ProductGroupWardrobeAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_WARDROBE_ACCESSORIES
	case daltymodel.ProductGroupInteriorDecoration:
		return proto.ProductGroup_PRODUCT_GROUP_INTERIOR_DECORATION
	case daltymodel.ProductGroupSleepTherapy:
		return proto.ProductGroup_PRODUCT_GROUP_SLEEP_THERAPY
	case daltymodel.ProductGroupKingKoil:
		return proto.ProductGroup_PRODUCT_GROUP_KING_KOIL
	case daltymodel.ProductGroupErgomotionAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_ERGOMOTION_ACCESSORIES
	case daltymodel.ProductGroupChildrenBedBases:
		return proto.ProductGroup_PRODUCT_GROUP_CHILDREN_BED_BASES
	case daltymodel.ProductGroupPillowCovers:
		return proto.ProductGroup_PRODUCT_GROUP_PILLOW_COVERS
	case daltymodel.ProductGroupTableware:
		return proto.ProductGroup_PRODUCT_GROUP_TABLEWARE
	case daltymodel.ProductGroupSets:
		return proto.ProductGroup_PRODUCT_GROUP_SETS
	case daltymodel.ProductGroupChildrenBedrooms:
		return proto.ProductGroup_PRODUCT_GROUP_CHILDREN_BEDROOMS
	case daltymodel.ProductGroupSpaceOrganizationStorage:
		return proto.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION_STORAGE
	case daltymodel.ProductGroupBathroomProducts:
		return proto.ProductGroup_PRODUCT_GROUP_BATHROOM_PRODUCTS
	case daltymodel.ProductGroupToys:
		return proto.ProductGroup_PRODUCT_GROUP_TOYS
	case daltymodel.ProductGroupAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_ACCESSORIES
	case daltymodel.ProductGroupNewYear:
		return proto.ProductGroup_PRODUCT_GROUP_NEW_YEAR
	case daltymodel.ProductGroupArmchairs:
		return proto.ProductGroup_PRODUCT_GROUP_ARMCHAIRS
	case daltymodel.ProductGroupMassageChairs:
		return proto.ProductGroup_PRODUCT_GROUP_MASSAGE_CHAIRS
	case daltymodel.ProductGroupKafkaTest:
		return proto.ProductGroup_PRODUCT_GROUP_KAFKA_TEST
	case daltymodel.ProductGroupCaseFurnitureAccessories:
		return proto.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE_ACCESSORIES
	default:
		return proto.ProductGroup_PRODUCT_GROUP_UNSPECIFIED
	}
}
