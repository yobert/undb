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
		var s;
		for(var path in this._grabs) {
			s = this._grabs[path];
			s.removeListener('CHANGE', this.storeUpdated);
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

		var s = this._grabs[path];
		if(!s) {
			s = this.store.Find(path);
			if(!s)
				throw("Grab failed: Unable to find path '" + path + "' in store '" + this.store.Name + "'");

			s.addListener('CHANGE', this.storeUpdated);
			this._grabs[path] = s;
		}

		return s;
	}
};

module.exports = mixin;
