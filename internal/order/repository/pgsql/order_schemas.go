package pgsql

var GetAllSuperAdminOrders = `SELECT * FROM public."orders" limit $1 offset $2`
var GetAllAdminOrders = `SELECT 
orders.id
orders.user_id
orders.address
orders.phone
orders.comment
orders.status
orders.create_ts
orders.update_ts
orders.active
orders.version
orders.notes
FROM public."orders"
INNER JOIN public."users" ON users.id = orders.user_id
WHERE users.id = $3 OR users.parent = $3
limit $1 offset $2`
var GetAllUserOrders = `SELECT 
orders.id,
orders.user_id,
orders.address,
orders.phone,
orders.comment,
orders.status,
orders.create_ts,
orders.update_ts,
orders.state,
orders.version,
orders.notes
FROM public."orders"
INNER JOIN public."users" ON users.id = $3
WHERE orders.user_id = $3 or users.role = 'ADMIN' 
limit $1 offset $2`

var GetUserOrdersByRole = `SELECT * FROM public."orders" limit $1 offset $2`

var GetAllAdminOrdersItem = `SELECT 
order_item.*,
sku.sku "sku.sku",
sku.small_name "sku.small_name" 
FROM public."order_item" 
INNER JOIN public."sku" ON order_item.sku_id = sku.id
WHERE order_item.order_id = $1`

var GetSingleOrder = `SELECT 
orders.id,
orders.user_id,
orders.address,
orders.phone,
orders.comment,
orders.status,
orders.notes
FROM public."orders"
WHERE orders.id = $1 AND orders.state = 'enabled';`

var CreateOrder = `INSERT 
INTO public."orders"(id, user_id, address, phone, comment, status, create_ts, update_ts, state, version)
VALUES (uuid_generate_v4(),$1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
`

var CreateOrderItem = `WITH I AS ( INSERT 
INTO public."order_item"(order_id, product_id, quantity, price, create_ts, update_ts, state, version)
SELECT $2, cart.sku_id, cart.quantity, sku.price,  $3, $4, $5, $6
FROM public."cart"
INNER JOIN public."users" ON users.id = $1
INNER JOIN public."sku" ON sku.id = cart.sku_id )
UPDATE public."cart"
SET state = 'deleted',
cart.update_ts = $4,
cart.version = version + 1
WHERE cart.user_id = $1;
`

var UpdateOrder = `UPDATE public."orders" AS o
SET
address = COALESCE(NULLIF($3, ''), o.address),
phone = COALESCE(NULLIF($4, ''), o.phone),
comment = COALESCE(NULLIF($5, ''), o.comment),
notes = COALESCE(NULLIF($6, ''), o.notes),
update_ts = COALESCE($7, o.update_ts),
version = o.version + 1
FROM public."orders"
INNER JOIN public."users" ON users.id = orders.user_id
WHERE o.id = $1 AND o.user_id = $2 OR users.role = 'ADMIN'
RETURNING o.id;`

var UpdateOrderStatus = `UPDATE public."orders" 
SET
status = $2,
update_ts =$3,
version = version + 1
WHERE orders.id = $1
RETURNING id;`

var DeleteOrder = `UPDATE public."orders" 
SET
state = 'deleted',
update_ts =$2,
version = version + 1
WHERE orders.id = $1
RETURNING id;`
