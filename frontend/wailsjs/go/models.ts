export namespace main {
	
	export class Log {
	    time: number;
	    type: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Log(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.type = source["type"];
	        this.message = source["message"];
	    }
	}
	export class Process {
	    id: number;
	    name: string;
	    create_time: number;
	    command: string;
	    status: number;
	    run_status: string;
	
	    static createFrom(source: any = {}) {
	        return new Process(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.create_time = source["create_time"];
	        this.command = source["command"];
	        this.status = source["status"];
	        this.run_status = source["run_status"];
	    }
	}

}

