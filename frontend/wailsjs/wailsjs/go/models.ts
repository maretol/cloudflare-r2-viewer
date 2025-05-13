export namespace request {
	
	export class BucketOption {
	    bucket_name: string;
	    description: string;
	    publish_url: string;
	
	    static createFrom(source: any = {}) {
	        return new BucketOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_name = source["bucket_name"];
	        this.description = source["description"];
	        this.publish_url = source["publish_url"];
	    }
	}

}

export namespace response {
	
	export class BucketList {
	    buckets: string[];
	
	    static createFrom(source: any = {}) {
	        return new BucketList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.buckets = source["buckets"];
	    }
	}
	export class BucketOption {
	    bucket_name: string;
	    description: string;
	    publish_url: string;
	
	    static createFrom(source: any = {}) {
	        return new BucketOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_name = source["bucket_name"];
	        this.description = source["description"];
	        this.publish_url = source["publish_url"];
	    }
	}
	export class Obj {
	    name: string;
	    path: string;
	    size: number;
	    publish_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Obj(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.publish_url = source["publish_url"];
	    }
	}
	export class ObjectList {
	    objects: Obj[];
	
	    static createFrom(source: any = {}) {
	        return new ObjectList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.objects = this.convertValues(source["objects"], Obj);
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

