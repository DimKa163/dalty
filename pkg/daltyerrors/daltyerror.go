package daltyerrors

import "fmt"

type DaltyErrorType int

const (
	DaltyErrorTypeWarning          DaltyErrorType = iota
	DaltyErrorTypeInvalidRequest                  = 3
	DaltyErrorTypeResourceNotFound                = 5
	DaltyErrorTypeBusinessError                   = 9
)

type DaltyErrorCode int

const (
	DaltyErrorCodeUnknown                                              DaltyErrorCode = iota
	DaltyErrorCodeRequiredFieldsMissing                                               = 1
	DaltyErrorCodeDeliveryCoordinatesMissing                                          = 2
	DaltyErrorCodePickupWarehouseMissing                                              = 3
	DaltyErrorCodeUnsupportedDeliveryMethod                                           = 4
	DaltyErrorCodeOutOfStock                                                          = 5
	DaltyErrorCodeResourceNotFound                                                    = 6
	DaltyErrorCodeDeliveryAddressNotSupported                                         = 7
	DaltyErrorCodeInsufficientFreeStock                                               = 8
	DaltyErrorCodeNoWarehouseForExpressDelivery                                       = 9
	DaltyErrorCodeInsufficientStockForPickup                                          = 10
	DaltyErrorCodeProductMissingProductionType                                        = 11
	DaltyErrorCodeOrderExecutionPlaceDeterminationFailed                              = 31
	DaltyErrorCodeMOLChainBuildingFailedNoCSCategory                                  = 32
	DaltyErrorCodeMOLChainBuildingFailedMOLDetermination                              = 33
	DaltyErrorCodeMOLChainBuildingFailedNearestWarehousesDetermination                = 35
	DaltyErrorCodeProductionSiteSelectionFailed                                       = 41
	DaltyErrorCodeProductionNormDeterminationFailed                                   = 42
	DaltyErrorCodeStopListExitDateDeterminationFailed                                 = 43
	DaltyErrorCodeDeprivationExitDateDeterminationFailed                              = 51
	DaltyErrorCodeProcurementLeadTimeCalculationFailed                                = 52
	DaltyErrorCodeDuplicateOrderItem                                                  = 53
	DaltyErrorCodeLogisticsLegToOrderExecutionPlaceBuildingFailed                     = 61
	DaltyErrorCodeLastMileCalculationFailed                                           = 71
	DaltyErrorCodeLastMileApihipCalculationFailed                                     = 72
)

type EntityError struct {
	ID         string
	EntityName string
}
type DaltyError struct {
	Code         DaltyErrorCode
	Type         DaltyErrorType
	Message      string
	Reason       string
	EntityErrors []*EntityError
}

func New(code DaltyErrorCode, msg string, entityErrors ...*EntityError) *DaltyError {
	if msg == "" {
		msg = "Ошибка при обработке запроса"
	}
	switch code {
	case DaltyErrorCodeRequiredFieldsMissing:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Не заполнены обязательные поля топика (продукт, количество)",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeDeliveryCoordinatesMissing:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "В одной или нескольких доставках установлен способ исполнения «Доставка», но не заданы координаты доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodePickupWarehouseMissing:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "В одной или нескольких доставках установлен способ исполнения «Самовывоз», но не выбран склад самовывоза",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeUnsupportedDeliveryMethod:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Указан неподдерживаемый сервисом способ доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeOutOfStock:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Одну или несколько позиций заказа невозможно исполнить из-за неактуальности МЦ или нахождения ее в дефиците и отсутствия ее на свободных остатках",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeResourceNotFound:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Не найден указанный филиал или продукт или склад",
			Type:         DaltyErrorTypeResourceNotFound,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeDeliveryAddressNotSupported:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Доставка по указанному адресу не производится",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case DaltyErrorCodeInsufficientFreeStock:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "В указаном месте обеспечения не хватает количества свободных остатков",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 9:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Не найден склад с остатками для экспресс-доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 10:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Не хватает остатков для самовывоза",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 11:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "У продукта не задан тип производства",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case 52:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Для отсутствующей продукции невозможно провести расчет срока ее закупки",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}
	case 53:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Дубль позиции заказа",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 31:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 32:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 33:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 35:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 41:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Ошибка выбора производственной площадки",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}

	case 42:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Ошибка определения нормы производства",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}

	case 43:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка определения нормы производства",
			EntityErrors: entityErrors,
		}

	case 51:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка определения нормы закупки",
			EntityErrors: entityErrors,
		}

	case 61:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка вычисления магистральной перевозки",
			EntityErrors: entityErrors,
		}

	case 71:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка расчета последней мили",
			EntityErrors: entityErrors,
		}

	case 72:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка расчета последней мили",
			EntityErrors: entityErrors,
		}
	default:
		return &DaltyError{
			Code:         code,
			Message:      msg,
			Reason:       "Не обрабатываемая ошибка",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}
	}
}

func (e *DaltyError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
