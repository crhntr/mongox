
db.runCommand({ collMod: "comments", validator: { $jsonSchema: { bsonType: "object", required: [ "_ac", "text"] } }, validationAction: "error" })
