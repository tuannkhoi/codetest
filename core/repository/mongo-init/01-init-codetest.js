const dbName = "codetest";
const dbHandle = db.getSiblingDB(dbName);

// Create the database and a default collection so it shows up immediately.
dbHandle.createCollection("event");
