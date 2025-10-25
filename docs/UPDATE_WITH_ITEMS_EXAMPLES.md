# Update with Items - Testing Examples

This document provides example request bodies for testing the PUT endpoints that update parent entities with all their items in a single transaction.

## Overview

All three endpoints follow the same pattern:
- **DELETE** all existing items (soft delete)
- **INSERT** new items from the request
- **UPDATE** parent entity
- All operations happen in a **single transaction**

## Authentication

All endpoints require JWT authentication. Include your JWT token in the request header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 1. Purchase Order - Update with Items

**Endpoint:** `PUT /api/v1/purchase-orders/:id/with-items`

**Description:** Updates a purchase order and replaces all its items.

### Request Body Example

```json
{
  "poNumber": "PO-2024-001-UPDATED",
  "userId": 1,
  "orderDate": "2024-01-15",
  "status": "approved",
  "totalAmount": 15000000,
  "description": "Updated purchase order with new items",
  "items": [
    {
      "productId": 1,
      "qtyOrdered": 50,
      "unitPrice": 150000,
      "totalPrice": 7500000,
      "description": "Updated product A"
    },
    {
      "productId": 2,
      "qtyOrdered": 30,
      "unitPrice": 250000,
      "totalPrice": 7500000,
      "description": "Updated product B"
    }
  ]
}
```

### Notes
- `poNumber`: Optional - if not provided, uses existing value
- `userId`: Optional - if 0 or not provided, uses existing value
- `orderDate`: Optional - if empty, uses existing value (format: YYYY-MM-DD)
- `status`: Optional - if empty, uses existing value (pending/approved/rejected/closed)
- `totalAmount`: Optional - if 0, uses existing value
- `description`: Optional - can be null
- `items`: **REQUIRED** - must have at least 1 item
- All existing items will be **soft deleted** and replaced with new items

### Full cURL Example

```bash
curl -X PUT http://localhost:8080/api/v1/purchase-orders/1/with-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "poNumber": "PO-2024-001-UPDATED",
    "userId": 1,
    "orderDate": "2024-01-15",
    "status": "approved",
    "totalAmount": 15000000,
    "description": "Updated purchase order",
    "items": [
      {
        "productId": 1,
        "qtyOrdered": 50,
        "unitPrice": 150000,
        "totalPrice": 7500000,
        "description": "Updated product A"
      },
      {
        "productId": 2,
        "qtyOrdered": 30,
        "unitPrice": 250000,
        "totalPrice": 7500000,
        "description": "Updated product B"
      }
    ]
  }'
```

---

## 2. Delivery Order - Update with Items

**Endpoint:** `PUT /api/v1/delivery-orders/:id/with-items`

**Description:** Updates a delivery order and replaces all its items.

### Request Body Example

```json
{
  "doNumber": "DO-2024-001-UPDATED",
  "purchaseOrderId": 1,
  "deliveryDate": "2024-01-20",
  "status": "shipped",
  "totalAmount": 12000000,
  "description": "Updated delivery order with modified items",
  "items": [
    {
      "productId": 1,
      "qtyDelivered": 45,
      "unitPrice": 150000,
      "totalPrice": 6750000,
      "description": "Delivered product A - updated quantity"
    },
    {
      "productId": 3,
      "qtyDelivered": 20,
      "unitPrice": 262500,
      "totalPrice": 5250000,
      "description": "New product C added to delivery"
    }
  ]
}
```

### Notes
- `doNumber`: Optional - if not provided, uses existing value
- `purchaseOrderId`: Optional - if 0 or not provided, uses existing value
- `deliveryDate`: Optional - if empty, uses existing value (format: YYYY-MM-DD)
- `status`: Optional - if empty, uses existing value (pending/in_transit/delivered/cancelled)
- `totalAmount`: Optional - if 0, uses existing value
- `description`: Optional - can be null
- `items`: **REQUIRED** - must have at least 1 item
- System validates that the purchase order exists

### Full cURL Example

```bash
curl -X PUT http://localhost:8080/api/v1/delivery-orders/1/with-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "doNumber": "DO-2024-001-UPDATED",
    "purchaseOrderId": 1,
    "deliveryDate": "2024-01-20",
    "status": "shipped",
    "totalAmount": 12000000,
    "description": "Updated delivery order",
    "items": [
      {
        "productId": 1,
        "qtyDelivered": 45,
        "unitPrice": 150000,
        "totalPrice": 6750000,
        "description": "Updated quantity"
      },
      {
        "productId": 3,
        "qtyDelivered": 20,
        "unitPrice": 262500,
        "totalPrice": 5250000,
        "description": "New product added"
      }
    ]
  }'
```

---

## 3. Invoice - Update with Items

**Endpoint:** `PUT /api/v1/invoices/:id/with-items`

**Description:** Updates an invoice and replaces all its items.

### Request Body Example

```json
{
  "invoiceNumber": "INV-2024-001-UPDATED",
  "userId": 2,
  "purchaseOrderId": 1,
  "deliveryOrderId": 1,
  "invoiceDate": "2024-01-25",
  "status": "sent",
  "totalAmount": 13500000,
  "description": "Updated invoice with corrected items",
  "items": [
    {
      "productId": 1,
      "qtyInvoiced": 45,
      "unitPrice": 150000,
      "totalPrice": 6750000,
      "description": "Invoice item A - updated"
    },
    {
      "productId": 2,
      "qtyInvoiced": 27,
      "unitPrice": 250000,
      "totalPrice": 6750000,
      "description": "Invoice item B - corrected quantity"
    }
  ]
}
```

### Notes
- `invoiceNumber`: Optional - if not provided, uses existing value
- `userId`: Optional - if 0 or not provided, uses existing value (reseller/customer ID)
- `purchaseOrderId`: Optional - can be null
- `deliveryOrderId`: Optional - can be null
- `invoiceDate`: Optional - if empty, uses existing value (format: YYYY-MM-DD)
- `status`: Optional - if empty, uses existing value (draft/sent/paid/closed)
- `totalAmount`: Optional - if 0, uses existing value
- `description`: Optional - can be null
- `items`: **REQUIRED** - must have at least 1 item
- System validates that user, PO (if provided), and DO (if provided) exist

### Full cURL Example

```bash
curl -X PUT http://localhost:8080/api/v1/invoices/1/with-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "invoiceNumber": "INV-2024-001-UPDATED",
    "userId": 2,
    "purchaseOrderId": 1,
    "deliveryOrderId": 1,
    "invoiceDate": "2024-01-25",
    "status": "sent",
    "totalAmount": 13500000,
    "description": "Updated invoice",
    "items": [
      {
        "productId": 1,
        "qtyInvoiced": 45,
        "unitPrice": 150000,
        "totalPrice": 6750000,
        "description": "Updated item A"
      },
      {
        "productId": 2,
        "qtyInvoiced": 27,
        "unitPrice": 250000,
        "totalPrice": 6750000,
        "description": "Updated item B"
      }
    ]
  }'
```

---

## Partial Update Examples

### Update Only Status and Items

You can update only specific fields. Other fields will keep their existing values.

**Purchase Order - Update Status and Items:**
```json
{
  "status": "approved",
  "items": [
    {
      "productId": 1,
      "qtyOrdered": 100,
      "unitPrice": 150000,
      "totalPrice": 15000000,
      "description": "Approved items"
    }
  ]
}
```

**Delivery Order - Update Status and Items:**
```json
{
  "status": "delivered",
  "items": [
    {
      "productId": 1,
      "qtyDelivered": 95,
      "unitPrice": 150000,
      "totalPrice": 14250000,
      "description": "Delivered items"
    }
  ]
}
```

**Invoice - Update Status and Items:**
```json
{
  "status": "paid",
  "items": [
    {
      "productId": 1,
      "qtyInvoiced": 95,
      "unitPrice": 150000,
      "totalPrice": 14250000,
      "description": "Paid items"
    }
  ]
}
```

---

## Response Format

### Success Response (200)

```json
{
  "success": true,
  "message": "Purchase Order with items updated successfully",
  "data": {
    "id": 1,
    "poNumber": "PO-2024-001-UPDATED",
    "userId": 1,
    "orderDate": "2024-01-15T00:00:00Z",
    "status": "approved",
    "totalAmount": 15000000,
    "description": "Updated purchase order",
    "user": {
      "id": 1,
      "name": "Admin User",
      "email": "admin@example.com"
    },
    "items": [
      {
        "id": 10,
        "purchaseOrderId": 1,
        "productId": 1,
        "qtyOrdered": 50,
        "unitPrice": 150000,
        "totalPrice": 7500000,
        "description": "Updated product A",
        "product": {
          "id": 1,
          "name": "Product A",
          "sku": "SKU-A001"
        }
      },
      {
        "id": 11,
        "purchaseOrderId": 1,
        "productId": 2,
        "qtyOrdered": 30,
        "unitPrice": 250000,
        "totalPrice": 7500000,
        "description": "Updated product B",
        "product": {
          "id": 2,
          "name": "Product B",
          "sku": "SKU-B002"
        }
      }
    ],
    "createdAt": "2024-01-10T10:00:00Z",
    "updatedAt": "2024-01-15T14:30:00Z"
  }
}
```

### Error Responses

**400 - Validation Error:**
```json
{
  "success": false,
  "message": "At least one item is required",
  "error": "Items array cannot be empty"
}
```

**404 - Not Found:**
```json
{
  "success": false,
  "message": "Purchase order not found",
  "error": "purchase order not found"
}
```

**409 - Conflict (Duplicate Number):**
```json
{
  "success": false,
  "message": "Purchase order number already exists",
  "error": "PO number 'PO-2024-001-UPDATED' already exists"
}
```

---

## Testing Workflow

### Recommended Testing Steps:

1. **Get existing data:**
   ```bash
   GET /api/v1/purchase-orders/1/with-items
   ```

2. **Update with modified items:**
   ```bash
   PUT /api/v1/purchase-orders/1/with-items
   # Use request body from examples above
   ```

3. **Verify the update:**
   ```bash
   GET /api/v1/purchase-orders/1/with-items
   # Check that old items are gone and new items exist
   ```

4. **Check soft delete (optional):**
   ```sql
   -- Old items should have deleted_at set
   SELECT * FROM purchase_order_items 
   WHERE purchase_order_id = 1 
   AND deleted_at IS NOT NULL;
   ```

---

## Important Notes

### Transaction Behavior
- All operations (update parent + delete old items + insert new items) happen in a **single transaction**
- If any step fails, the entire transaction is **rolled back**
- No partial updates will occur

### Soft Delete Pattern
- Old items are **soft deleted** (deleted_at timestamp is set)
- They remain in the database for audit trail purposes
- They will not appear in normal queries (GORM automatically filters them out)

### Validation Rules
1. **Items Required:** Must provide at least 1 item in the items array
2. **Foreign Key Validation:** 
   - User ID must exist
   - Product IDs in items must exist
   - Purchase Order ID must exist (for Delivery Order and Invoice)
   - Delivery Order ID must exist (for Invoice, if provided)
3. **Unique Constraints:**
   - PO Number must be unique (or same as current)
   - DO Number must be unique (or same as current)
   - Invoice Number must be unique (or same as current)
4. **Date Format:** All dates must be in YYYY-MM-DD format

### Status Values
- **Purchase Order:** pending, approved, rejected, closed
- **Delivery Order:** pending, in_transit, delivered, cancelled
- **Invoice:** draft, sent, paid, closed

---

## Postman Collection

You can import these examples into Postman:

1. Create a new collection
2. Add environment variables:
   - `base_url`: http://localhost:8080
   - `jwt_token`: Your authentication token
3. Create requests for each endpoint using the examples above
4. Set Authorization header: `Bearer {{jwt_token}}`

---

## Common Issues & Solutions

### Issue: "At least one item is required"
**Solution:** Ensure the `items` array has at least one item.

### Issue: "Invalid date format"
**Solution:** Use YYYY-MM-DD format (e.g., "2024-01-15").

### Issue: "User not authenticated"
**Solution:** Include valid JWT token in Authorization header.

### Issue: "Purchase order number already exists"
**Solution:** Either use a different number or omit the field to keep existing number.

### Issue: "Product not found"
**Solution:** Verify all productIds in items array exist in the database.

---

## Related Endpoints

- **GET with items:** `/api/v1/purchase-orders/:id/with-items`
- **POST with items:** `/api/v1/purchase-orders/with-items`
- **Filtered GET:** `/api/v1/purchase-orders/filter?user_id=1&status=approved`

Same pattern applies for delivery-orders and invoices.
