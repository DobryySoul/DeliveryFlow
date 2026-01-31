package nats

import "slices"

const (
	SubjectOrderCreated      = "events.order_created"
	SubjectInventoryReserved = "events.inventory_reserved"
	SubjectPaymentCaptured   = "events.payment_captured"
	SubjectDeliveryAssigned  = "events.delivery_assigned"
	SubjectNotifyUser        = "jobs.notify_user"
	SubjectCreateOrder       = "rpc.create_order"
	SubjectGetOrderStatus    = "rpc.get_order_status"
)

var Subjects = []string{
	SubjectOrderCreated,
	SubjectInventoryReserved,
	SubjectPaymentCaptured,
	SubjectDeliveryAssigned,
	SubjectNotifyUser,
	SubjectCreateOrder,
	SubjectGetOrderStatus,
}

func IsValidSubject(subject string) bool {
	return slices.Contains(Subjects, subject)
}
