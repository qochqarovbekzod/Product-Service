CURRENT_DIR=$(shell pwd)

DBURL := postgres://postgres:1918@localhost:5432/products?sslmode=disable


proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}


swag-init:
	swag init -g api/routes.go --output api/docs


mig-up:
	migrate -path databases/migrations -database '${DBURL}' -verbose up

mig-down:
	migrate -path databases/migrations -database '${DBURL}' -verbose down

mig-force:
	migrate -path databases/migrations -database '${DBURL}' -verbose force 1


mig-create-product:
	migrate create -ext sql -dir databases/migrations -seq create_product_table

mig-create-orders:
	migrate create -ext sql -dir databases/migrations -seq create_orders_table

mig-create-orderItems:
	migrate create -ext sql -dir databases/migrations -seq create_order_items_table

mig-create-productcategories:
	migrate create -ext sql -dir databases/migrations -seq create_product_categories_table

mig-create-reting:
	migrate create -ext sql -dir databases/migrations -seq create_reting_table

mig-create-payment:
	migrate create -ext sql -dir databases/migrations -seq create_payment_table




