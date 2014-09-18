"use strict";

var events = require('events');

// todo: use polyfill for IE < 9 support for object.keys()

var STORES = 1;
var VALUES = 2;

var INSERT = 1;
var DELETE = 2;
var UPDATE = 3;
var MERGE = 4;

function Store(id, typ) {
	this.Id = id;
	this.Type = typ;
	this.Records = {};
	this.Deleted = false;

	if(typ == STORES)
		this.next = 0;

	return;
}

Store.prototype.addListener = events.EventEmitter.prototype.addListener;
Store.prototype.removeListener = events.EventEmitter.prototype.removeListener;
Store.prototype.emit = events.EventEmitter.prototype.emit;

Store.prototype.Insert = function(insert, source) {
	if(insert.Id in this.Records)
		return "Insert into store '" + this.Id + "' failed: '" + insert.Id + "' already exists";

	insert['parent'] = this;
	this.Records[insert.Id] = insert;
	this.emitChange({Method: INSERT, Id: insert.Id, Type: insert.Type}, source);
	return undefined;
}

Store.prototype.Delete = function(source) {
	this.Deleted = true;
	this.emitChange({Method: DELETE}, source);
	return undefined;
}

Store.prototype.Update = function(values, source) {
	if(this.Type != VALUES)
		return "Update store '" + this.Id + "' failed: not a VALUES store";

	this.Records = values;
	this.emitChange({Method: UPDATE, Values: values}, source);
	return undefined;
}

Store.prototype.Merge = function(values, source) {
	if(this.Type != VALUES)
		return "Merge store '" + this.Id + "' failed: not a VALUES store";

	for(var k in values)
		this.Records[k] = values[k];

	this.emitChange({Method: MERGE, Values: values}, source);
	return undefined;
}

Store.prototype.Path = function() {
	var p = this.Id;
	var pt = this['parent'];
	while(pt) {
		p = pt.Id + '.' + p;
		pt = pt['parent'];
	}
	return p;
}

Store.prototype.Find = function(path) {
	var chunks = path.split(/\./);
	if(!chunks.length || chunks[0] != this.Id) {
		return undefined;
	}

	var s = this;
	for(var i = 1; i < chunks.length; i++) {
		if(s.Type != STORES)
			return undefined;

		s = s.Records[chunks[i]]
		if(!s)
			return undefined;
	}

	return s
}

Store.prototype.FindOrThrow = function(path) {
	var s = this.Find(path);

	if(!s)
		throw("Find on store '" + this.Id + "' failed: Find path '" + path + "' failed");

	return s;
}

Store.prototype.Seq = function() {
	this.next++;
	// sloppy hack for the moment
	return('jsid_' + new Date().getTime() + this.next + Math.floor(Math.random() * 100000));
}

Store.prototype.Exec = function(op, changesource) {
	var s = this.Find(op.Path);
	if(!s)
		return "Exec on store '" + this.Id + "' failed: Find path '" + op.Path + "' failed";

	switch(op.Method) {
	case INSERT:
		return s.Insert(new Store(op.Id, op.Type), changesource);
	case DELETE:
		return s.Delete(changesource);
	case UPDATE:
		return s.Update(op.Values, changesource);
	case MERGE:
		return s.Merge(op.Values, changesource);
	}
	return "Exec on store '" + this.Id + "' failed: Invalid method: '" + op.Method + "'";
}

Store.prototype.emitChange = function(op, changesource) {
	var s = this;
	var path = s.Id;

	while(s) {
		s.emit('CHANGE', {
			Method: op.Method,
			Path: path,
			Id: op.Id,
			Type: op.Type,
			Values: op.Values,
			changesource: changesource
		});
		s = s['parent'];
		if(s) {
			path = s.Id + '.' + path;
		}
	}
}

module.exports = {
	STORES: STORES,
	VALUES: VALUES,

	INSERT: INSERT,
	DELETE: DELETE,
	UPDATE: UPDATE,
	MERGE: MERGE,

	Store: Store
};

