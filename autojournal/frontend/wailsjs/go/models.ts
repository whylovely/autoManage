export namespace domain {
	
	export class Backup {
	    ID: number;
	    FilePath: string;
	    Note: string;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Backup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FilePath = source["FilePath"];
	        this.Note = source["Note"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
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
	export class Expense {
	    ID: number;
	    VehicleID: number;
	    CategoryID: number;
	    Amount: number;
	    OdometerAt: number;
	    // Go type: time
	    Date: any;
	    Description: string;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Expense(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.VehicleID = source["VehicleID"];
	        this.CategoryID = source["CategoryID"];
	        this.Amount = source["Amount"];
	        this.OdometerAt = source["OdometerAt"];
	        this.Date = this.convertValues(source["Date"], null);
	        this.Description = source["Description"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
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
	export class Vehicle {
	    id: number;
	    vin: string;
	    make: string;
	    Model: string;
	    Year: number;
	    color?: string;
	    EngineVolume: number;
	    FuelType: number;
	    Odometer: number;
	    Notes: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdateAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Vehicle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vin = source["vin"];
	        this.make = source["make"];
	        this.Model = source["Model"];
	        this.Year = source["Year"];
	        this.color = source["color"];
	        this.EngineVolume = source["EngineVolume"];
	        this.FuelType = source["FuelType"];
	        this.Odometer = source["Odometer"];
	        this.Notes = source["Notes"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdateAt = this.convertValues(source["UpdateAt"], null);
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

