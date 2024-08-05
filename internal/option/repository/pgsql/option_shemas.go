package pgsql

var GetSkuValueProducts = `SELECT
sku_value.id,
sku_value.option_id, 
sku_value.option_value_id
FROM public."sku_value"
WHERE sku_value.sku_id = $1
`

var GetProductOptions = `SELECT
option.id, 
option.name
FROM public."option"
WHERE option.id = $1 AND option.state = 'enabled'`

var GetProductOptionsByCat = `SELECT
option.id, 
option.name
FROM public."option"
WHERE option.category_id=$1  AND option.state = 'enabled'`

var GetOptionValues = `SELECT 
option_value.id,
option_value.name
FROM public."option_value"
WHERE option_value.id = $1 AND option_value.state = 'enabled'`

var GetOptionValuesByOptId = `SELECT 
option_value.id,
option_value.name
FROM public."option_value"
WHERE option_value.option_id = $1  AND option_value.state = 'enabled'`

var GetOptionBySkuValueId = `SELECT 
option.id,
option.name,
option_value.id,
option_value.name
FROM public."sku_value"
INNER JOIN public."option" ON sku_value.option_id = option.id 
INNER JOIN public."option_value" ON sku_value.option_value_id = option_value.id
WHERE sku_value.id = $1 AND sku_value.state = 'enabled' AND option.state = 'enabled';`

var CreateOption = `INSERT 
INTO public."option"(category_id, name, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`
var CreateOptionValue = `INSERT 
INTO public."option_value"(option_id, name)
VALUES ($1, $2)
RETURNING id;
`
var CreateSkuValue = `INSERT 
INTO public."sku_value"(sku_id, option_id, option_value_id)
VALUES ($1, $2, $3)
RETURNING id;
`

var UpdateOptionName = `UPDATE public."option"
SET name = $2,
category_id = $3,
update_ts =$4,
version = version + 1
WHERE option.id = $1;
`

var UpdateOptionValueName = `UPDATE public."option_value"
SET name = $2,
option_id = $3,
update_ts =$4,
version = version + 1
WHERE option_value.id = $1;
`

var RemoveOption = `UPDATE public."option"
SET state = 'deleted',
version = version + 1
WHERE option.id = $1;
`
var RemoveOptionValue = `UPDATE public."option_value"
SET state = 'deleted',
version = version + 1
WHERE option_value.id = $1;
`
var RemoveProductSkuValue = `DELETE FROM public."sku_value"
WHERE sku_value.id = $1;
`
