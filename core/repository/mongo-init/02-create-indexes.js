const dbName = "codetest";
const dbHandle = db.getSiblingDB(dbName);

// Indexes for search filters.
dbHandle.event.createIndex({ "bettingstatus.value": 1 });
dbHandle.event.createIndex({ "eventvisibility.value": 1 });
dbHandle.event.createIndex({ "startTimeBSONDate": 1 });
