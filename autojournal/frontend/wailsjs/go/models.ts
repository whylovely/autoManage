export namespace domain {
	
	export class Backup {
	    id: number;
	    filePath: string;
	    note: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Backup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.filePath = source["filePath"];
	        this.note = source["note"];
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
	export class Reminder {
	    id: number;
	    vehicleId: number;
	    title: string;
	    reminderType: string;
	    intervalKm?: number;
	    intervalDays?: number;
	    lastDoneOdometer?: number;
	    // Go type: time
	    lastDoneDate?: any;
	    // Go type: time
	    nextDueDate?: any;
	    nextDueOdometer?: number;
	    isActive: boolean;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Reminder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vehicleId = source["vehicleId"];
	        this.title = source["title"];
	        this.reminderType = source["reminderType"];
	        this.intervalKm = source["intervalKm"];
	        this.intervalDays = source["intervalDays"];
	        this.lastDoneOdometer = source["lastDoneOdometer"];
	        this.lastDoneDate = this.convertValues(source["lastDoneDate"], null);
	        this.nextDueDate = this.convertValues(source["nextDueDate"], null);
	        this.nextDueOdometer = source["nextDueOdometer"];
	        this.isActive = source["isActive"];
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
	export class DueReminder {
	    reminder: Reminder;
	    dueByDate: boolean;
	    dueByOdometer: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DueReminder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reminder = this.convertValues(source["reminder"], Reminder);
	        this.dueByDate = source["dueByDate"];
	        this.dueByOdometer = source["dueByOdometer"];
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
	    id: number;
	    vehicleId: number;
	    categoryId: number;
	    amount: number;
	    odometerAt: number;
	    // Go type: time
	    date: any;
	    description: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Expense(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vehicleId = source["vehicleId"];
	        this.categoryId = source["categoryId"];
	        this.amount = source["amount"];
	        this.odometerAt = source["odometerAt"];
	        this.date = this.convertValues(source["date"], null);
	        this.description = source["description"];
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
	export class ExpenseCategory {
	    id: number;
	    name: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	    }
	}
	export class ExpenseCategoryTotal {
	    categoryId: number;
	    categoryName: string;
	    totalAmount: number;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseCategoryTotal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.totalAmount = source["totalAmount"];
	    }
	}
	export class MonthlyExpenseTotal {
	    month: string;
	    totalAmount: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyExpenseTotal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.totalAmount = source["totalAmount"];
	    }
	}
	export class ExpenseStats {
	    totalAmount: number;
	    byCategory: ExpenseCategoryTotal[];
	    byMonth: MonthlyExpenseTotal[];
	
	    static createFrom(source: any = {}) {
	        return new ExpenseStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalAmount = source["totalAmount"];
	        this.byCategory = this.convertValues(source["byCategory"], ExpenseCategoryTotal);
	        this.byMonth = this.convertValues(source["byMonth"], MonthlyExpenseTotal);
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
	    model: string;
	    year: number;
	    color?: string;
	    engineVolume: number;
	    fuelType: number;
	    odometer: number;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Vehicle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vin = source["vin"];
	        this.make = source["make"];
	        this.model = source["model"];
	        this.year = source["year"];
	        this.color = source["color"];
	        this.engineVolume = source["engineVolume"];
	        this.fuelType = source["fuelType"];
	        this.odometer = source["odometer"];
	        this.notes = source["notes"];
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

