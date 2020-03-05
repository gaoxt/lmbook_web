const DB_NAME = 'bookdb';
const DB_VERSION = 1;
let DB;

var idba = {
	async getDb() {
		return new Promise((resolve, reject) => {

			if (DB) { return resolve(DB); }
			console.log('OPENING DB', DB);
			let request = window.indexedDB.open(DB_NAME, DB_VERSION);

			request.onerror = e => {
				console.log('Error opening db', e);
				reject('Error');
			};

			request.onsuccess = e => {
				DB = e.target.result;
				resolve(DB);
			};

			request.onupgradeneeded = e => {
				console.log('onupgradeneeded');
				let db = e.target.result;
				db.createObjectStore("books", { keyPath: 'id' });
				db.createObjectStore("expire", { keyPath: 'id' });
				db.createObjectStore("bookDetail", { keyPath: 'id' });
			};
		});
	},
	async getBooks() {
		let db = await this.getDb();
		return new Promise(resolve => {
			let trans = db.transaction(['books'], 'readonly');
			trans.oncomplete = () => {
				resolve(books);
			};
			let store = trans.objectStore('books');
			let books = [];

			store.openCursor().onsuccess = e => {
				let cursor = e.target.result;
				if (cursor) {
					books.push(cursor.value)
					cursor.continue();
				}
			};

		});
	},

	async getBookItem(id) {
		let db = await this.getDb();
		return new Promise(resolve => {
			let trans = db.transaction(['books'], 'readonly');
			trans.oncomplete = () => {
				resolve(bookItem);
			};
			let store = trans.objectStore('books');
			let getter = store.get(id);
			let bookItem = [];
			getter.onsuccess = e => {
				bookItem = e.target.result;
			};
		});
	},

	async saveBook(book) {
		let db = await this.getDb();
		return new Promise(resolve => {
			let trans = db.transaction(['books'], 'readwrite');
			trans.oncomplete = () => {
				resolve();
			};
			let store = trans.objectStore('books');
			store.put(book);
		});
	},


	async getBookDetail(id) {
		let db = await this.getDb();
		return new Promise(resolve => {
			let trans = db.transaction(['bookDetail'], 'readonly');
			trans.oncomplete = () => {
				resolve(bookDetail);
			};
			let store = trans.objectStore('bookDetail');
			let bookDetail = [];

			let getter = store.get(id);
			getter.onsuccess = e => {
				bookDetail = e.target.result.bookDetail;
			};
		});
	},

	async saveBookDetail(id, bookDetail) {
		let db = await this.getDb();
		return new Promise((resolve) => {
			let trans = db.transaction(['bookDetail'], 'readwrite');
			trans.oncomplete = () => {
				resolve(bookDetail);
			};
			let store = trans.objectStore('bookDetail');
			store.put({ id: id, bookDetail });
		});
	},

	async saveExpire(key) {
		let db = await this.getDb();
		return new Promise(resolve => {
			let trans = db.transaction(['expire'], 'readwrite');
			trans.oncomplete = () => {
				resolve();
			};
			let store = trans.objectStore('expire');
			store.put({
				id: key,
				expire: Date.now(),
			});
		});
	},
	// 默认1小时过期
	async isExpire(key, expireTime = 60 * 60 * 1) {
		let db = await this.getDb();
		return new Promise((resolve, reject) => {
			let trans = db.transaction(['expire'], 'readonly');
			let store = trans.objectStore('expire');
			let getter = store.get(key);
			getter.onsuccess = e => {
				if (e.target.result && ((new Date().getTime() - e.target.result.expire) < expireTime * 1000)) {
					resolve();
				}
				reject(new Error('timeout!'))
			};
			getter.onerror = e => {
				reject(e)
			};
		});
	}
}



var idb = {
	state: {
		bookList: [],
		bookDetail: {},
		isExpire: false,
	},
	mutations: {
		setExpire(state, isExpire) {
			state.isExpire = isExpire;
		},
		setBookList(state, data) {
			state.bookList = data;
		},
		setBookDetail(state, data) {
			state.bookDetail = data;
		},
	},
	actions: {
		async getBooks({ commit }) {
			let books = await idba.getBooks();
			commit('setBookList', books.reverse());
		},
		async saveBook(context, bookList) {
			Object.values(bookList).reverse().forEach((c) => {
				let bookInfo = {};
				Object.keys(c).forEach(function (key) {
					if (key === "Id") {
						bookInfo.id = Number(c[key]);
					} else {
						bookInfo[key] = c[key];
					}
				});
				idba.saveBook(bookInfo);
			});
			await idba.saveExpire("books");
		},
		async getBookDetail({ commit }, id) {
			let detail = await idba.getBookDetail(id);
			let main = await idba.getBookItem(Number(id));
			main.Details = detail
			commit('setBookDetail', main);
		},
		async saveBookDetail(context, data) {
			await idba.saveBookDetail(data.id, data.res);
			await idba.saveExpire("detail_" + data.id);
		},
		async checkExpire({ commit }, expireKey) {
			try {
				await idba.isExpire(expireKey, 60 * 60 * 24 * 1);
				commit('setExpire', false);
			} catch (e) {
				await idba.saveExpire(expireKey);
				commit('setExpire', true);
			}
		}
	}
}

export default { idb };