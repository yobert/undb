// this uses object.keys
// which is isn't in IE < 9

function Init(name) {
	return {
		Name: name,
		Records: {},
		deleted: false
	}
}

function Insert(store, name, record) {
	if(name in store.Records)
		return "Insert into store '" + store.Name + "' failed: '" + name + "' already exists";

	store.Records[name] = record
	return undefined;
}

function Update(store, name, record) {
	if(!(name in store.Records)) {
		return "Update store '" + store.Name + "' failed: '" + name + "' does not exist";
	}
	store.Records[name] = record;
	return undefined;
}

function Upsert(store, name, record) {
	store.Records[name] = record;
	return undefined;
}

function Delete(store) {
	store.Deleted = true;
	return undefined;
}

function Get(store, name) {
	var exists = false;
	if(name in store.Records)
		exists = true;

	return [store.Records, exists];
}

function Keys(store) {
	return store.Records.keys();
}

function Exec(store, op) {
	switch(op.Method) {
	case "Insert":
		return [undefined, Insert(store, op.Name, op.Record)];
	case "Update":
		return [undefined, Update(store, op.Name, op.Record)];
	case "Upsert":
		Upsert(store, op.Name, op.Record)
		return [undefined, undefined];
	case "Delete":
		Delete(store)
		return [undefined, undefined];
	case "Get":
		var ve = Get(store, op.Name);
		if !ve[1] {
			return [undefined, undefined];
		}
		return [ve[0], undefined];
	case "Keys":
		return [Keys(store), undefined];
	}

	return [undefined, "invalid method: '" + op.Method + "'"];
}

module.exports = {
	Init: Init,
	Insert: Insert,
	Update: Update,
	Upsert: Upsert,
	Delete: Delete,
	Get: Get,
	Keys: Keys,
	Exec: Exec
}
