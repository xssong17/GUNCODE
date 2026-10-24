export namespace main {
	
	export class AppInfo {
	    datasetVersion: string;
	    databasePath: string;
	    weaponCount: number;
	    categoryCount: number;
	    loadoutCount: number;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.datasetVersion = source["datasetVersion"];
	        this.databasePath = source["databasePath"];
	        this.weaponCount = source["weaponCount"];
	        this.categoryCount = source["categoryCount"];
	        this.loadoutCount = source["loadoutCount"];
	    }
	}

}

export namespace model {
	
	export class Category {
	    id: number;
	    key: string;
	    name: string;
	    order: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.key = source["key"];
	        this.name = source["name"];
	        this.order = source["order"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportResult {
	    imported: number;
	    skipped: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imported = source["imported"];
	        this.skipped = source["skipped"];
	    }
	}
	export class LoadoutDTO {
	    id: number;
	    weaponId: number;
	    weaponKey: string;
	    weaponName: string;
	    weaponInGameName: string;
	    categoryKey: string;
	    categoryName: string;
	    title: string;
	    mode: string;
	    code: string;
	    note: string;
	    tuning: string;
	    tags: string[];
	    favorite: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new LoadoutDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.weaponId = source["weaponId"];
	        this.weaponKey = source["weaponKey"];
	        this.weaponName = source["weaponName"];
	        this.weaponInGameName = source["weaponInGameName"];
	        this.categoryKey = source["categoryKey"];
	        this.categoryName = source["categoryName"];
	        this.title = source["title"];
	        this.mode = source["mode"];
	        this.code = source["code"];
	        this.note = source["note"];
	        this.tuning = source["tuning"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LoadoutInput {
	    weaponId: number;
	    title: string;
	    mode: string;
	    code: string;
	    note: string;
	    tuning: string;
	    tags: string[];
	    favorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LoadoutInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.weaponId = source["weaponId"];
	        this.title = source["title"];
	        this.mode = source["mode"];
	        this.code = source["code"];
	        this.note = source["note"];
	        this.tuning = source["tuning"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	    }
	}
	export class ParsedCode {
	    raw: string;
	    matched: boolean;
	    weaponKey: string;
	    weaponName: string;
	    weaponInGameName: string;
	    categoryKey: string;
	    categoryName: string;
	    mode: string;
	    payload: string;
	
	    static createFrom(source: any = {}) {
	        return new ParsedCode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.raw = source["raw"];
	        this.matched = source["matched"];
	        this.weaponKey = source["weaponKey"];
	        this.weaponName = source["weaponName"];
	        this.weaponInGameName = source["weaponInGameName"];
	        this.categoryKey = source["categoryKey"];
	        this.categoryName = source["categoryName"];
	        this.mode = source["mode"];
	        this.payload = source["payload"];
	    }
	}
	export class Weapon {
	    id: number;
	    key: string;
	    name: string;
	    inGameName: string;
	    categoryId: number;
	    categoryKey: string;
	    categoryName: string;
	    aliases: string;
	    caliber: string;
	    status: string;
	    source: string;
	    note: string;
	    isCustom: boolean;
	    sort: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Weapon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.key = source["key"];
	        this.name = source["name"];
	        this.inGameName = source["inGameName"];
	        this.categoryId = source["categoryId"];
	        this.categoryKey = source["categoryKey"];
	        this.categoryName = source["categoryName"];
	        this.aliases = source["aliases"];
	        this.caliber = source["caliber"];
	        this.status = source["status"];
	        this.source = source["source"];
	        this.note = source["note"];
	        this.isCustom = source["isCustom"];
	        this.sort = source["sort"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

