package daltymodel

import "github.com/beevik/guid"

type (
	PickupStrategy int
	Line           struct {
		Product  *Product       `json:"product"`
		Quantity int32          `json:"quantity"`
		Strategy PickupStrategy `json:"strategy"`
	}
	SpecificationType int
	Specification     struct {
		Product       *Line             `json:"product"`
		Type          SpecificationType `json:"type"`
		Strategy      PickupStrategy    `json:"strategy"`
		ChildProducts []*Line           `json:"child_products"`
	}
	SpecResponse struct {
		Specifications []*Specification `json:"specifications"`
	}
)

const (
	SpecificationTypeDefault SpecificationType = iota
	SpecificationTypeDirect
	SpecificationTypeReverse
)

const (
	PickupStrategyNearest PickupStrategy = iota
	PickupStrategyFarthest
)

type ProductType string

const (
	ProductTypeSKU           ProductType = "7ad7cb12-e3be-4a24-ad60-2c15dff8152b"
	ProductTypeMaterialAsset             = "e183d6ef-f33b-4dbe-b8cc-e0c176e4e78a"
)

func (p ProductType) String() string {
	switch p {
	case ProductTypeSKU:
		return "SKU"
	case ProductTypeMaterialAsset:
		return "MaterialAsset"
	default:
		return "Unknown"
	}
}

type ProductionType string

const (
	ProductionTypeUnknown    ProductionType = "UNKNOWN"
	ProductionTypeProducing                 = "f950f1c5-c871-4da0-b374-aa6eb1aa24d1"
	ProductionTypePurchasing                = "2cb8891a-f707-4bdc-a492-aef60f03f451"
)

func (p ProductionType) String() string {
	switch p {
	case ProductionTypeProducing:
		return "PRODUCTION"
	case ProductionTypePurchasing:
		return "PURCHASING"
	default:
		return string(ProductionTypeUnknown)
	}
}

type ProductGroup string

const (
	ProductGroupUnknown                    ProductGroup = "Unknown"
	ProductGroupKitchens                   ProductGroup = "14b85d2e-da37-4a50-8244-0616a0b31794"
	ProductGroupCaseFurniture              ProductGroup = "21644026-d9ca-4da2-b91e-1adb2507a824"
	ProductGroupBeddingSets                ProductGroup = "6d31ec4d-33a5-48f0-b9a2-3666a9fab68b"
	ProductGroupSofas                      ProductGroup = "f09f1c26-12f2-4d50-8c94-3f347b2970a2"
	ProductGroupCovers                     ProductGroup = "c8a49a04-07c6-4538-9054-478abcca71e5"
	ProductGroupBlankets                   ProductGroup = "66d86f8e-6812-4182-82cc-4792937e299a"
	ProductGroupBedBasesWithStorage        ProductGroup = "2743d9f1-1244-440f-addb-653efe23492c"
	ProductGroupSofaComponents             ProductGroup = "1362ba98-6ea2-45a4-9134-67813a17d2a5"
	ProductGroupErgomotion                 ProductGroup = "5471ed6b-c477-4853-a718-6b6bdfacb5ad"
	ProductGroupNonProducts                ProductGroup = "7c6ad459-836a-4796-87bd-77f5dab8cac1"
	ProductGroupSmallFurniture             ProductGroup = "260d757a-bc45-4e95-9225-96cdf0cfe92f"
	ProductGroupMattresses                 ProductGroup = "81025548-4083-40bb-a0e3-a22652c500b5"
	ProductGroupSlattedBases               ProductGroup = "66f0375b-9789-49e1-81f9-a824d649e3c8"
	ProductGroupMattressToppers            ProductGroup = "8ac5ad9c-1772-4c37-86be-ab7c35c64761"
	ProductGroupPillows                    ProductGroup = "868b44a7-40a3-4152-bc75-b41e597d1d98"
	ProductGroupBeds                       ProductGroup = "c7a26159-5856-49df-aed2-d053f672ae11"
	ProductGroupBedBases                   ProductGroup = "bd9c804c-a698-4075-b8bf-dcc197e33f0a"
	ProductGroupMiscellaneous              ProductGroup = "117d995d-9ae7-487f-99fe-eb0c3fe5d460"
	ProductGroupCaseFurnitureAccessories   ProductGroup = "c3dae860-b36c-4133-8407-8f7e4280f65d"
	ProductGroupMurphyBeds                 ProductGroup = "37d0e4b6-54ac-47ca-9a9a-83389118fb0d"
	ProductGroupWardrobes                  ProductGroup = "d717a93d-ab66-4e57-95c3-98c9b7a8c8df"
	ProductGroupBedAccessories             ProductGroup = "6d4aee05-19a2-4f8f-9d0b-2a6b7c41e189"
	ProductGroupMurphyBedAccessories       ProductGroup = "a286d6cf-31e7-494a-b17c-7d1d9a87d6d1"
	ProductGroupSmallFurnitureAccessories  ProductGroup = "51c5e588-ef43-496d-80b5-5d5832619cae"
	ProductGroupWardrobeAccessories        ProductGroup = "5ee02ed2-01a0-43c8-bb97-7675f01fbc64"
	ProductGroupInteriorDecoration         ProductGroup = "ecb374a7-3eaf-4e55-a381-2ce38ee64c02"
	ProductGroupTextiles                   ProductGroup = "6dc4e692-4f2a-4aed-bb91-d3aa14cf8b5b"
	ProductGroupSleepTherapy               ProductGroup = "a41f277e-378c-4b7e-a16f-fa08a01a40cf"
	ProductGroupElectronics                ProductGroup = "accfdbbb-628f-4611-bda2-9f912c6f26bd"
	ProductGroupClothing                   ProductGroup = "de5b8449-a610-427c-a345-f24664789a9d"
	ProductGroupOrthopedics                ProductGroup = "e9b7cd10-5145-4e28-8491-269adb048c6a"
	ProductGroupCoffeeTables               ProductGroup = "671afa27-9707-4cc2-9719-21be9b8281fb"
	ProductGroupKingKoil                   ProductGroup = "ef530893-fcb2-420d-85a1-3ea798495db9"
	ProductGroupErgomotionAccessories      ProductGroup = "37d04227-ddb1-4e79-91e6-e781b769689c"
	ProductGroupChildrenBedBases           ProductGroup = "c1073c8b-866e-455f-8170-0114466de42f"
	ProductGroupPillowCovers               ProductGroup = "83b31fe3-03b9-499a-8cb3-ddd60188516e"
	ProductGroupTableware                  ProductGroup = "f0383727-77da-410b-bf57-d62e850d87f2"
	ProductGroupSets                       ProductGroup = "729ff321-96e4-4bec-9e7e-eec9573ae58d"
	ProductGroupHomeOffice                 ProductGroup = "2aefbe7a-0c26-4e0c-84ae-9d3098528560"
	ProductGroupChildrenBedrooms           ProductGroup = "08537964-33d6-454e-b422-247896afc313"
	ProductGroupSpaceOrganizationStorage   ProductGroup = "9bc21671-bf65-4570-bdaf-344af5c22157"
	ProductGroupBathroomProducts           ProductGroup = "d570fea1-a978-4915-8de0-d4aa0aa1ac79"
	ProductGroupToys                       ProductGroup = "64f47b87-5359-473d-946c-a540287fce78"
	ProductGroupAccessories                ProductGroup = "856b85f1-2639-440f-b904-af0c579943d1"
	ProductGroupNewYear                    ProductGroup = "926f2f98-4de7-4c80-856d-57bd6c04847d"
	ProductGroupArmchairs                  ProductGroup = "fadf4ffa-55f9-4674-b254-ec6fafd67be4"
	ProductGroupMassageChairs              ProductGroup = "d390de83-2a44-417a-873b-485ecce030ca"
	ProductGroupLivingRooms                ProductGroup = "f8a304ef-212d-419e-82ed-61292bb0db31"
	ProductGroupLighting                   ProductGroup = "9f611676-2fd6-4852-be00-a17796592612"
	ProductGroupKafkaTest                  ProductGroup = "aec329db-d9e4-42db-af89-cb0ae05efc00"
	ProductGroupDecor                      ProductGroup = "c8dd79ea-b42b-4cc5-ad74-1d7862c4d253"
	ProductGroupSpaceOrganizationUpper     ProductGroup = "e65c8693-421b-4e50-97b7-df79ae23789c"
	ProductGroupHomeCareUpper              ProductGroup = "1c90ac25-9c16-48f7-bc8b-8a809e7c2705"
	ProductGroupSpaceOrganization          ProductGroup = "bb2724ce-aa43-4c27-9e54-a2d3247d68cc"
	ProductGroupHomeCare                   ProductGroup = "6a139e59-5578-4c66-8a30-8abf844ec40c"
	ProductGroupHallways                   ProductGroup = "7c6ed323-06a6-431a-b549-616f4f43bdf5"
	ProductGroupFurnitureProtectionAndCare ProductGroup = "99bb0a71-37bd-496d-b261-a8ffd1adcb8b"
	ProductGroupOutdoorFurniture           ProductGroup = "4ce1132d-03d7-4cb4-84d7-325c5751aca7"
	ProductGroupStorage                    ProductGroup = "2b360d0c-b5a4-4db8-b682-ead00d0c2c58"
	ProductGroupInterior                   ProductGroup = "c1357f61-440c-4f1e-aeb4-15a26895e886"
	ProductGroupSeasonalProducts           ProductGroup = "4d7ec861-a530-4c5a-b133-12d173fca3be"
	ProductGroupFragrances                 ProductGroup = "189995ff-cb88-408a-9ea4-c82672f71b22"
)

type Product struct {
	ID                  guid.Guid      `json:"id"`
	Name                string         `json:"name"`
	Type                ProductType    `json:"type"`
	Fnrec               string         `json:"fnrec"`
	IsArchive           bool           `json:"is_archive"`
	IntegrationID       string         `json:"integration_id"`
	IsService           bool           `json:"is_service"`
	ProductionType      ProductionType `json:"production_type"`
	Group               ProductGroup   `json:"group"`
	SeriesID            string         `json:"series_id"`
	CategoryID          string         `json:"category_id"`
	AccountProviderId   string         `json:"account_provider_id"`
	NonStandardCategory string         `json:"non_standard_category"`
	CountMa             int32          `json:"count_ma"`
	Volume              float64        `json:"volume"`
	Length              float64        `json:"length"`
	Width               float64        `json:"width"`
	Height              float64        `json:"height"`
	Weight              float64        `json:"weight"`
}

func NewLine(product *Product, quantity int32, strategy PickupStrategy) *Line {
	return &Line{
		Product:  product,
		Quantity: quantity,
		Strategy: strategy,
	}
}

func NewDefaultSpecification(line *Line, strategy PickupStrategy) *Specification {
	return &Specification{
		Product:       line,
		Type:          SpecificationTypeDefault,
		Strategy:      strategy,
		ChildProducts: make([]*Line, 0),
	}
}

func NewDirectSpecification(line *Line, strategy PickupStrategy, subProducts []*Line) *Specification {
	return &Specification{
		Product:       line,
		Type:          SpecificationTypeDirect,
		Strategy:      strategy,
		ChildProducts: subProducts,
	}
}

func NewReverseSpecification(line *Line, strategy PickupStrategy, subProducts []*Line) *Specification {
	return &Specification{
		Product:       line,
		Type:          SpecificationTypeReverse,
		Strategy:      strategy,
		ChildProducts: subProducts,
	}
}
