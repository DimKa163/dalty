package daltyerrors

import "fmt"

type DaltyErrorType int

const (
	DaltyErrorTypeWarning          DaltyErrorType = iota
	DaltyErrorTypeInvalidRequest                  = 3
	DaltyErrorTypeResourceNotFound                = 5
	DaltyErrorTypeBusinessError                   = 9
)

type EntityError struct {
	ID         string
	EntityName string
}
type DaltyError struct {
	Code         int
	Type         DaltyErrorType
	Message      string
	Reason       string
	EntityErrors []*EntityError
}

func New(code int, entityErrors ...*EntityError) *DaltyError {
	switch code {
	case 1:
		return &DaltyError{
			Code:         code,
			Message:      "Не заполнены обязательные поля топика (продукт, количество)",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 2:
		return &DaltyError{
			Code:         code,
			Message:      "В одной или нескольких доставках установлен способ исполнения «Доставка», но не заданы координаты доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 3:
		return &DaltyError{
			Code:         code,
			Message:      "В одной или нескольких доставках установлен способ исполнения «Самовывоз», но не выбран склад самовывоза",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 4:
		return &DaltyError{
			Code:         code,
			Message:      "Указан неподдерживаемый сервисом способ доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 5:
		return &DaltyError{
			Code:         code,
			Message:      "Одну или несколько позиций заказа невозможно исполнить из-за неактуальности МЦ или нахождения ее в дефиците и отсутствия ее на свободных остатках",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case 6:
		return &DaltyError{
			Code:         code,
			Message:      "Не найден указанный филиал или продукт или склад",
			Type:         DaltyErrorTypeResourceNotFound,
			EntityErrors: entityErrors,
		}

	case 7:
		return &DaltyError{
			Code:         code,
			Message:      "Доставка по указанному адресу не производится",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case 8:
		return &DaltyError{
			Code:         code,
			Message:      "В указаном месте обеспечения не хватает количества свободных остатков",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 9:
		return &DaltyError{
			Code:         code,
			Message:      "Не найден склад с остатками для экспресс-доставки",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 10:
		return &DaltyError{
			Code:         code,
			Message:      "Не хватает остатков для самовывоза",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 11:
		return &DaltyError{
			Code:         code,
			Message:      "У продукта не задан тип производства",
			Type:         DaltyErrorTypeBusinessError,
			EntityErrors: entityErrors,
		}

	case 52:
		return &DaltyError{
			Code:         code,
			Message:      "Для отсутствующей продукции невозможно провести расчет срока ее закупки",
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка расчета закупки",
			EntityErrors: entityErrors,
		}

	case 53:
		return &DaltyError{
			Code:         code,
			Message:      "Дубль позиции заказа",
			Type:         DaltyErrorTypeInvalidRequest,
			EntityErrors: entityErrors,
		}

	case 31:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения места исполнения заказа",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 32:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка построения цепочки МОЛ (нет категории МОЛ=\"ЦС\")",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 33:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения МОЛ для сбора свободных остатков",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 35:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения перечня ближайших складов обеспечения",
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка построения цепочки МОЛов",
			EntityErrors: entityErrors,
		}

	case 41:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка выбора производственной площадки",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}

	case 42:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения нормы производства",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}

	case 43:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения срока выхода из стоп-листа",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка определения нормы производства",
			EntityErrors: entityErrors,
		}

	case 51:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка определения срока выхода из дефицита",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка определения нормы закупки",
			EntityErrors: entityErrors,
		}

	case 61:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка построения логистического плеча к месту исполнения заказа",
			Type:         DaltyErrorTypeWarning,
			Reason:       "Ошибка вычисления магистральной перевозки",
			EntityErrors: entityErrors,
		}

	case 71:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка построения плеча с места обеспечения",
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка расчета последней мили",
			EntityErrors: entityErrors,
		}

	case 72:
		return &DaltyError{
			Code:         code,
			Message:      "Ошибка расчета срока последней мили Apihip",
			Type:         DaltyErrorTypeBusinessError,
			Reason:       "Ошибка расчета последней мили",
			EntityErrors: entityErrors,
		}
	default:
		return &DaltyError{
			Code:         code,
			Message:      "Не обрабатываемая ошибка",
			Type:         DaltyErrorTypeWarning,
			EntityErrors: entityErrors,
		}
	}
}

func (e *DaltyError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
