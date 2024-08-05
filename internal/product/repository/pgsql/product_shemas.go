package pgsql

var GetSkuProducts = `SELECT 
sku.*,
product.product_name "product.product_name",
product.description "product.description",
product.category_id "product.category_id",
product.language "product.language",
product.create_ts "product.create_ts",
FROM public."sku"
INNER JOIN public."product" ON sku.product_id = product.id 
limit $1 offset $2;`

var GetSkuSingleProduct = `SELECT 
sku.id,
sku.product_id,
sku.sku,
sku.price,
sku.quantity,
sku.large_name,
sku.small_name,
sku.thumb_name,
sku.count_viewed,
sku.create_ts,
sku.update_ts,
sku.state,
sku.version,
product.product_name "product.product_name",
product.description "product.description",
product.category_id "product.category_id",
product.create_ts "product.create_ts"
FROM public."sku"
INNER JOIN public."product" ON sku.product_id = product.id 
WHERE sku.sku = $1;`

var GetSkuProductsByCategory = `SELECT 
sku.id,
sku.product_id,
sku.sku,
sku.price,
sku.quantity,
sku.large_name,
sku.small_name,
sku.thumb_name,
sku.count_viewed,
sku.create_ts,
sku.update_ts,
sku.state,
sku.version,
product.product_name "product.product_name",
product.description "product.description",
product.category_id "product.category_id",
product.create_ts "product.create_ts",
count(*) OVER() AS total_count
FROM public."sku"
INNER JOIN public."product" ON sku.product_id = product.id 
WHERE product.category_id = $3 AND  product.state = 'enabled'
limit $1 offset $2;`

var ProductsbyCategory = `SELECT
product.id "product.id",
product.product_name "product.product_name",
product.description "product.description",
product.category_id "product.category_id",
product.create_ts "product.create_ts",
count(*) OVER() AS total_count
FROM public."product"
WHERE product.category_id = $3 AND product.state = 'enabled'
limit $1 offset $2;`

var SkuByProdID = `SELECT
sku.id "sku.id",
sku.sku "sku.sku",
sku.price "sku.price",
sku.quantity "sku.quantity",
sku.small_name "sku.small_name"
FROM public."sku"
WHERE sku.product_id = $1`

var GetSkuValueProducts = `SELECT
sku_value.id,
sku_value.option_id, 
sku_value.option_value_id
FROM public."sku_value"
Where sku_value.skuId = $1`

var GetProductOptions = `SELECT
option.id, 
option.name
FROM public."option"
Where option.id = $1`

var GetOptionValues = `SELECT 
option_value.id,
option_value.name
FROM public."option_value"
Where option_value.id = $1`

var CreateProduct = `INSERT 
INTO public."product"(product_name, description, category_id, brand_id, region_id, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
`

var CreateSku = `INSERT 
INTO public."sku"(product_id, sku, price, quantity, large_name, small_name, thumb_name, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11);
`
var CreateOption = `INSERT 
INTO public."option"(product_id, name, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`
var CreateOptionValue = `INSERT 
INTO public."option_value"(product_id, option_id, name, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
`
var CreateSkuValue = `INSERT 
INTO public."sku_value"(skuId, product_id, option_id, option_value_id, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;
`

var UpdateSku = `UPDATE public."sku"
SET price = $2,
	quantity = $3,
	large_name = $4,
	small_name = $5,
	thumb_name = $6,
	create_ts = create_ts,
	update_ts = $7,
	state = $8,
	version = version + 1
WHERE sku.id = $1;
`

var UpdateProduct = `UPDATE public."product"
SET category_id = $2,
	product_name = $3,
	description = $4,
	brand_id = $5,
	region_id = $6,
	update_ts = $7,
	state = $8,
	version = version + 1
WHERE product.id = $1;
`

var UpdateOptionValueName = `UPDATE public."option_value"
SET name = $2,
version = version + 1
WHERE option_value.id = $1;
`

var Removeproduct = `UPDATE public."product"
SET state = 'deleted',
version = version + 1
WHERE product.id = $1;
`

var RemoveSku = `UPDATE public."sku"
SET state = 'deleted',
version = version + 1
WHERE sku.id = $1;
`

var RemoveProductOption = `DELETE FROM public."option"
WHERE option.id = $1;
`
var RemoveProductOptionValue = `DELETE FROM public."option_value"
WHERE option_value.id = $1;
`
var RemoveProductSkuValue = `DELETE FROM public."sku_value"
WHERE sku_value.id = $1;
`
