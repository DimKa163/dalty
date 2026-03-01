package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		var daltyErr *daltyerrors.DaltyError
		if errors.As(err, &daltyErr) {
			return nil, protoerr.Handle(daltyErr)
		}
		return nil, status.Error(codes.Internal, err.Error())
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
	request.Specs = make(map[string]*usecase.Spec)

	for _, line := range lines {
		if !line.HasId() {
			err = fmt.Errorf("invalid specification line: %s", line)
		}
		if !line.HasQuantity() {
			err = fmt.Errorf("invalid specification line: %w", err)
		}
		if !line.HasFnrec() && !line.HasIntegration() {
			err = fmt.Errorf("invalid specification line: %w", err)
		}
		request.Specs[line.GetId()] = &usecase.Spec{
			ID:            line.GetId(),
			IntegrationID: line.GetIntegration(),
			Fnrec:         line.GetFnrec(),
			Quantity:      line.GetQuantity(),
			RelateToID:    line.GetRelateToId(),
		}
	}
	return
}

func toSpecification(spec *daltymodel.Specification) *proto.Specification {
	var specification proto.Specification
	if spec.Product != nil {
		specification.SetProduct(toLine(spec.Product))
	}
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
	line.SetId(ln.ID)
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

func toProtoProductionType(in daltymodel.ProductionType) proto.ProductionType {
	switch in {
	case daltymodel.ProductionTypeProducing:
		return proto.ProductionType_PRODUCTION_TYPE_PRODUCING
	case daltymodel.ProductionTypePurchasing:
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
