// A helper mixin for react components.
// Use by setting this.store to your root store,
// and call this.grab(path) to get a store link that
// will also trigger forceUpdate of this component
// if the data changes.
//
// TODO: figure out react event batching

var mixin = {
	componentWillMount: function() {
		this._grabs = {};
	},

	componentWillUnmount: function() {
		var s, storename, path;
		for(storename in this._grabs) {
			for(path in this._grabs[storename]) {
				s = this._grabs[storename][path];
				s.removeListener('CHANGE', this.storeUpdated);
			}
		}
		delete this._grabs;
	},

	storeUpdated: function() {
		if(this.isMounted()) {
			this.forceUpdate(this.storeUpdateFinished);
		}
	},

	grab: function(path) {
		if(!this.store)
			throw('grab() called before this.store set');

		return this.grabFrom(this.store, path);
	},

	grabFrom: function(store, path) {
		var storename = store.Id;
		var gs = this._grabs[storename];
		if(!gs) {
			gs = {};
			this._grabs[storename] = gs;
		}

		var s = gs[path];
		if(!s) {
			s = store.Find(path);
			if(!s)
				throw("Grab failed: Unable to find path '" + path + "' in store '" + store.Id + "'");

			s.addListener('CHANGE', this.storeUpdated);
			gs[path] = s;
		}

		return s;
	}
};

module.exports = mixin;
