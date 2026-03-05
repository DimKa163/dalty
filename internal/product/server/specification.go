package server

import (
	"context"
	"errors"
	"fmt"
	productsv1 "github.com/DimKa163/dalty/pkg/api/products/v1"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DimKa163/dalty/internal/product/usecase"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"google.golang.org/grpc"
)

type SpecificationServer struct {
	app *usecase.SpecificationService
	productsv1.UnimplementedSpecificationServiceServer
}

func NewSpecificationServer(app *usecase.SpecificationService) *SpecificationServer {
	return &SpecificationServer{app: app}
}

func (ss *SpecificationServer) Bind(server *grpc.Server) {
	productsv1.RegisterSpecificationServiceServer(server, ss)
}

func (ss *SpecificationServer) Execute(ctx context.Context, in *productsv1.SpecificationRequest) (*productsv1.SpecificationResponse, error) {
	var response productsv1.SpecificationResponse
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
	specs := make([]*productsv1.Specification, len(res))
	for i, r := range res {
		specs[i] = toSpecification(r)
	}
	response.SetSpecifications(specs)
	return &response, nil
}

func toIn(in *productsv1.SpecificationRequest, request *usecase.SpecRequest) (err error) {
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

func toSpecification(spec *daltymodel.Specification) *productsv1.Specification {
	var specification productsv1.Specification
	if spec.Product != nil {
		specification.SetProduct(toLine(spec.Product))
	}
	specification.SetType(productsv1.SpecificationType(spec.Type))
	specification.SetStrategy(productsv1.PickupStrategy(spec.Strategy))
	childProducts := make([]*productsv1.Line, len(spec.ChildProducts))
	for i, childProduct := range spec.ChildProducts {
		childProducts[i] = toLine(childProduct)
	}
	specification.SetChildProduct(childProducts)
	return &specification
}

func toLine(ln *daltymodel.Line) *productsv1.Line {
	var line productsv1.Line
	line.SetId(ln.ID)
	line.SetProduct(toproductsv1ProductV2(ln.Product))
	line.SetQuantity(ln.Quantity)
	line.SetStrategy(productsv1.PickupStrategy(ln.Strategy))
	return &line
}

func toproductsv1ProductV2(in *daltymodel.Product) *productsv1.Product {
	var out productsv1.Product
	out.SetId(in.ID.String())
	out.SetName(in.Name)
	out.SetProductionType(toproductsv1ProductionType(in.ProductionType))
	out.SetFnrec(in.Fnrec)
	out.SetIsService(in.IsService)
	out.SetGroup(toproductsv1ProductGroup(in.Group))
	out.SetSeriesId(in.SeriesID)
	out.SetCategoryId(in.CategoryID)
	out.SetAccountProvider(in.AccountProviderId)
	out.SetNonStandardCategoryId(in.NonStandardCategory)
	var pack productsv1.Pack
	pack.SetLength(in.Length)
	pack.SetHeight(in.Height)
	pack.SetWidth(in.Width)
	pack.SetVolume(in.Volume)
	pack.SetWeight(in.Weight)
	out.SetPack(&pack)
	return &out
}

func toproductsv1ProductionType(in daltymodel.ProductionType) productsv1.ProductionType {
	switch in {
	case daltymodel.ProductionTypeProducing:
		return productsv1.ProductionType_PRODUCTION_TYPE_PRODUCING
	case daltymodel.ProductionTypePurchasing:
		return productsv1.ProductionType_PRODUCTION_TYPE_PURCHASING
	default:
		return productsv1.ProductionType_PRODUCTION_TYPE_UNKNOWN
	}
}

func toproductsv1ProductGroup(pg daltymodel.ProductGroup) productsv1.ProductGroup {
	switch pg {
	case daltymodel.ProductGroupKitchens:
		return productsv1.ProductGroup_PRODUCT_GROUP_KITCHENS
	case daltymodel.ProductGroupCaseFurniture:
		return productsv1.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE
	case daltymodel.ProductGroupBeddingSets:
		return productsv1.ProductGroup_PRODUCT_GROUP_BEDDING_SETS
	case daltymodel.ProductGroupSofas:
		return productsv1.ProductGroup_PRODUCT_GROUP_SOFAS
	case daltymodel.ProductGroupCovers:
		return productsv1.ProductGroup_PRODUCT_GROUP_COVERS
	case daltymodel.ProductGroupBlankets:
		return productsv1.ProductGroup_PRODUCT_GROUP_BLANKETS
	case daltymodel.ProductGroupBedBasesWithStorage:
		return productsv1.ProductGroup_PRODUCT_GROUP_BED_BASES_WITH_STORAGE
	case daltymodel.ProductGroupSofaComponents:
		return productsv1.ProductGroup_PRODUCT_GROUP_SOFA_COMPONENTS
	case daltymodel.ProductGroupErgomotion:
		return productsv1.ProductGroup_PRODUCT_GROUP_ERGOMOTION
	case daltymodel.ProductGroupNonProducts:
		return productsv1.ProductGroup_PRODUCT_GROUP_NON_PRODUCTS
	case daltymodel.ProductGroupSmallFurniture:
		return productsv1.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE
	case daltymodel.ProductGroupMattresses:
		return productsv1.ProductGroup_PRODUCT_GROUP_MATTRESSES
	case daltymodel.ProductGroupSlattedBases:
		return productsv1.ProductGroup_PRODUCT_GROUP_SLATTED_BASES
	case daltymodel.ProductGroupMattressToppers:
		return productsv1.ProductGroup_PRODUCT_GROUP_MATTRESS_TOPPERS
	case daltymodel.ProductGroupPillows:
		return productsv1.ProductGroup_PRODUCT_GROUP_PILLOWS
	case daltymodel.ProductGroupBeds:
		return productsv1.ProductGroup_PRODUCT_GROUP_BEDS
	case daltymodel.ProductGroupBedBases:
		return productsv1.ProductGroup_PRODUCT_GROUP_BED_BASES
	case daltymodel.ProductGroupMiscellaneous:
		return productsv1.ProductGroup_PRODUCT_GROUP_MISCELLANEOUS
	case daltymodel.ProductGroupWardrobes:
		return productsv1.ProductGroup_PRODUCT_GROUP_WARDROBES
	case daltymodel.ProductGroupTextiles:
		return productsv1.ProductGroup_PRODUCT_GROUP_TEXTILES
	case daltymodel.ProductGroupElectronics:
		return productsv1.ProductGroup_PRODUCT_GROUP_ELECTRONICS
	case daltymodel.ProductGroupClothing:
		return productsv1.ProductGroup_PRODUCT_GROUP_CLOTHING
	case daltymodel.ProductGroupOrthopedics:
		return productsv1.ProductGroup_PRODUCT_GROUP_ORTHOPEDICS
	case daltymodel.ProductGroupCoffeeTables:
		return productsv1.ProductGroup_PRODUCT_GROUP_COFFEE_TABLES
	case daltymodel.ProductGroupHomeOffice:
		return productsv1.ProductGroup_PRODUCT_GROUP_HOME_OFFICE
	case daltymodel.ProductGroupLivingRooms:
		return productsv1.ProductGroup_PRODUCT_GROUP_LIVING_ROOMS
	case daltymodel.ProductGroupLighting:
		return productsv1.ProductGroup_PRODUCT_GROUP_LIGHTING
	case daltymodel.ProductGroupDecor:
		return productsv1.ProductGroup_PRODUCT_GROUP_DECOR
	case daltymodel.ProductGroupSpaceOrganizationUpper:
		return productsv1.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case daltymodel.ProductGroupHomeCareUpper:
		return productsv1.ProductGroup_PRODUCT_GROUP_HOME_CARE_UPPER
	case daltymodel.ProductGroupSpaceOrganization:
		return productsv1.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION
	case daltymodel.ProductGroupHomeCare:
		return productsv1.ProductGroup_PRODUCT_GROUP_HOME_CARE
	case daltymodel.ProductGroupHallways:
		return productsv1.ProductGroup_PRODUCT_GROUP_HALLWAYS
	case daltymodel.ProductGroupFurnitureProtectionAndCare:
		return productsv1.ProductGroup_PRODUCT_GROUP_FURNITURE_PROTECTION_AND_CARE
	case daltymodel.ProductGroupOutdoorFurniture:
		return productsv1.ProductGroup_PRODUCT_GROUP_OUTDOOR_FURNITURE
	case daltymodel.ProductGroupStorage:
		return productsv1.ProductGroup_PRODUCT_GROUP_STORAGE
	case daltymodel.ProductGroupInterior:
		return productsv1.ProductGroup_PRODUCT_GROUP_INTERIOR
	case daltymodel.ProductGroupSeasonalProducts:
		return productsv1.ProductGroup_PRODUCT_GROUP_SEASONAL_PRODUCTS
	case daltymodel.ProductGroupFragrances:
		return productsv1.ProductGroup_PRODUCT_GROUP_FRAGRANCES

	// accessories & additional groups
	case daltymodel.ProductGroupMurphyBeds:
		return productsv1.ProductGroup_PRODUCT_GROUP_MURPHY_BEDS
	case daltymodel.ProductGroupBedAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_BED_ACCESSORIES
	case daltymodel.ProductGroupMurphyBedAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_MURPHY_BED_ACCESSORIES
	case daltymodel.ProductGroupSmallFurnitureAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_SMALL_FURNITURE_ACCESSORIES
	case daltymodel.ProductGroupWardrobeAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_WARDROBE_ACCESSORIES
	case daltymodel.ProductGroupInteriorDecoration:
		return productsv1.ProductGroup_PRODUCT_GROUP_INTERIOR_DECORATION
	case daltymodel.ProductGroupSleepTherapy:
		return productsv1.ProductGroup_PRODUCT_GROUP_SLEEP_THERAPY
	case daltymodel.ProductGroupKingKoil:
		return productsv1.ProductGroup_PRODUCT_GROUP_KING_KOIL
	case daltymodel.ProductGroupErgomotionAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_ERGOMOTION_ACCESSORIES
	case daltymodel.ProductGroupChildrenBedBases:
		return productsv1.ProductGroup_PRODUCT_GROUP_CHILDREN_BED_BASES
	case daltymodel.ProductGroupPillowCovers:
		return productsv1.ProductGroup_PRODUCT_GROUP_PILLOW_COVERS
	case daltymodel.ProductGroupTableware:
		return productsv1.ProductGroup_PRODUCT_GROUP_TABLEWARE
	case daltymodel.ProductGroupSets:
		return productsv1.ProductGroup_PRODUCT_GROUP_SETS
	case daltymodel.ProductGroupChildrenBedrooms:
		return productsv1.ProductGroup_PRODUCT_GROUP_CHILDREN_BEDROOMS
	case daltymodel.ProductGroupSpaceOrganizationStorage:
		return productsv1.ProductGroup_PRODUCT_GROUP_SPACE_ORGANIZATION_STORAGE
	case daltymodel.ProductGroupBathroomProducts:
		return productsv1.ProductGroup_PRODUCT_GROUP_BATHROOM_PRODUCTS
	case daltymodel.ProductGroupToys:
		return productsv1.ProductGroup_PRODUCT_GROUP_TOYS
	case daltymodel.ProductGroupAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_ACCESSORIES
	case daltymodel.ProductGroupNewYear:
		return productsv1.ProductGroup_PRODUCT_GROUP_NEW_YEAR
	case daltymodel.ProductGroupArmchairs:
		return productsv1.ProductGroup_PRODUCT_GROUP_ARMCHAIRS
	case daltymodel.ProductGroupMassageChairs:
		return productsv1.ProductGroup_PRODUCT_GROUP_MASSAGE_CHAIRS
	case daltymodel.ProductGroupKafkaTest:
		return productsv1.ProductGroup_PRODUCT_GROUP_KAFKA_TEST
	case daltymodel.ProductGroupCaseFurnitureAccessories:
		return productsv1.ProductGroup_PRODUCT_GROUP_CASE_FURNITURE_ACCESSORIES
	default:
		return productsv1.ProductGroup_PRODUCT_GROUP_UNSPECIFIED
	}
}
