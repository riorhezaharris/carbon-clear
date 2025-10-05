// MongoDB initialization script for order service
db = db.getSiblingDB('carbon_clear_orders');

// Create collections with validation
db.createCollection('orders', {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["user_id", "project_id", "quantity", "status", "created_at"],
      properties: {
        user_id: {
          bsonType: "string",
          description: "User ID is required and must be a string"
        },
        project_id: {
          bsonType: "string",
          description: "Project ID is required and must be a string"
        },
        quantity: {
          bsonType: "number",
          minimum: 1,
          description: "Quantity must be a positive number"
        },
        status: {
          bsonType: "string",
          enum: ["pending", "confirmed", "completed", "cancelled"],
          description: "Status must be one of the enum values"
        },
        created_at: {
          bsonType: "date",
          description: "Created at must be a date"
        },
        updated_at: {
          bsonType: "date",
          description: "Updated at must be a date"
        }
      }
    }
  }
});

db.createCollection('carts', {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["user_id", "items", "created_at"],
      properties: {
        user_id: {
          bsonType: "string",
          description: "User ID is required and must be a string"
        },
        items: {
          bsonType: "array",
          description: "Items must be an array"
        },
        created_at: {
          bsonType: "date",
          description: "Created at must be a date"
        },
        updated_at: {
          bsonType: "date",
          description: "Updated at must be a date"
        }
      }
    }
  }
});

// Create indexes for better performance
db.orders.createIndex({ "user_id": 1 });
db.orders.createIndex({ "project_id": 1 });
db.orders.createIndex({ "status": 1 });
db.orders.createIndex({ "created_at": 1 });

db.carts.createIndex({ "user_id": 1 });
db.carts.createIndex({ "created_at": 1 });

print('MongoDB initialization completed for order service');
