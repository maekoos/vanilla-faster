-- name: GetProductsForSubcategory :many
SELECT * FROM products
WHERE subcategory_slug=$1
ORDER BY slug asc;


-- name: GetCollections :many
SELECT * FROM collections
ORDER BY collections.name;

-- name: GetCollectionsWithCategories :many
SELECT sqlc.embed(collections), sqlc.embed(categories) FROM collections
LEFT JOIN categories ON categories.collection_id=collections.id
ORDER BY collections.name, categories.name;


-- name: GetProductDetails :one
SELECT * FROM products
WHERE slug=$1;

-- name: GetSubcategory :one
SELECT * FROM subcategories
WHERE slug=$1;


-- name: GetCategory :many
SELECT sqlc.embed(categories), sqlc.embed(subcollections), sqlc.embed(subcategories) FROM categories
LEFT JOIN subcollections ON subcollections.category_slug=categories.slug
LEFT JOIN subcategories ON subcategories.subcollection_id=subcollections.id
WHERE categories.slug=$1;


-- name: GetCollectionDetails :many
SELECT sqlc.embed(collections), sqlc.embed(categories) FROM collections
LEFT JOIN categories ON categories.collection_id = collections.id
WHERE collections.slug=$1
ORDER BY collections.slug asc;
-- with: {
-- categories: true,
-- },


-- name: GetProductCount :one
SELECT count(*) as total FROM products;


-- name: GetCategoryProductCount :one
SELECT count(*) FROM categories
LEFT JOIN subcollections ON categories.slug = subcollections.category_slug
LEFT JOIN subcategories ON subcollections.id = subcategories.subcollection_id
LEFT JOIN products ON subcategories.slug = products.subcategory_slug
WHERE categories.slug=$1;


-- name: GetSubcategoryProductCount :one
SELECT count(*) FROM products
WHERE products.subcategory_slug=$1;