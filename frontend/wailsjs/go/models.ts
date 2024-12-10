export namespace contract {
	
	export class Review {
	    Body: string;
	    State: string;
	    Author: string;
	
	    static createFrom(source: any = {}) {
	        return new Review(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Body = source["Body"];
	        this.State = source["State"];
	        this.Author = source["Author"];
	    }
	}
	export class Pullrequest {
	    Number: number;
	    URL: string;
	    Author: string;
	    Title: string;
	    // Go type: time
	    Opened: any;
	    Assignee: string;
	    Reviewers: string[];
	    Description: string;
	    Mergeable: boolean;
	    Draft: boolean;
	    Reviews: Review[];
	
	    static createFrom(source: any = {}) {
	        return new Pullrequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Number = source["Number"];
	        this.URL = source["URL"];
	        this.Author = source["Author"];
	        this.Title = source["Title"];
	        this.Opened = this.convertValues(source["Opened"], null);
	        this.Assignee = source["Assignee"];
	        this.Reviewers = source["Reviewers"];
	        this.Description = source["Description"];
	        this.Mergeable = source["Mergeable"];
	        this.Draft = source["Draft"];
	        this.Reviews = this.convertValues(source["Reviews"], Review);
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
	export class Repository {
	    Name: string;
	    Owner: string;
	    Provider: string;
	
	    static createFrom(source: any = {}) {
	        return new Repository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Owner = source["Owner"];
	        this.Provider = source["Provider"];
	    }
	}

}

