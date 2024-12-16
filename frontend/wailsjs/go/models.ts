export namespace database {
	
	export class IPCount {
	    IP: string;
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new IPCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IP = source["IP"];
	        this.Count = source["Count"];
	    }
	}
	export class StatusCount {
	    Status: number;
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new StatusCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.Count = source["Count"];
	    }
	}

}

export namespace parser {
	
	export class LogLine {
	    IP: string;
	    Timestamp: string;
	    Method: string;
	    Path: string;
	    HttpVesion: string;
	    Status: number;
	    Size: number;
	    Filename: string;
	    LineNumber: number;
	
	    static createFrom(source: any = {}) {
	        return new LogLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IP = source["IP"];
	        this.Timestamp = source["Timestamp"];
	        this.Method = source["Method"];
	        this.Path = source["Path"];
	        this.HttpVesion = source["HttpVesion"];
	        this.Status = source["Status"];
	        this.Size = source["Size"];
	        this.Filename = source["Filename"];
	        this.LineNumber = source["LineNumber"];
	    }
	}

}

export namespace resultcollector {
	
	export class Result {
	    scope: string;
	    success: boolean;
	    filename: string;
	    line: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scope = source["scope"];
	        this.success = source["success"];
	        this.filename = source["filename"];
	        this.line = source["line"];
	        this.message = source["message"];
	    }
	}

}

