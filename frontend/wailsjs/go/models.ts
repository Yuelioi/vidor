export namespace app {
	
	export class App {
	
	
	    static createFrom(source: any = {}) {
	        return new App(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Part {
	
	
	    static createFrom(source: any = {}) {
	        return new Part(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class taskMap {
	
	
	    static createFrom(source: any = {}) {
	        return new taskMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace config {
	
	export class Config {
	    system?: models.SystemConfig;
	    plugins: {[key: string]: models.PluginConfig};
	    test?: models.SystemConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.system = this.convertValues(source["system"], models.SystemConfig);
	        this.plugins = this.convertValues(source["plugins"], models.PluginConfig, true);
	        this.test = this.convertValues(source["test"], models.SystemConfig);
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

export namespace models {
	
	export class PluginConfig {
	    id: string;
	    enable: boolean;
	    settings: {[key: string]: string};
	
	    static createFrom(source: any = {}) {
	        return new PluginConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.enable = source["enable"];
	        this.settings = source["settings"];
	    }
	}
	export class SystemConfig {
	    theme: string;
	    scale_factor: number;
	    proxy_url: string;
	    use_proxy: boolean;
	    magic_name: string;
	    download_dir: string;
	    download_video: boolean;
	    download_audio: boolean;
	    download_subtitle: boolean;
	    download_combine: boolean;
	    download_limit: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.scale_factor = source["scale_factor"];
	        this.proxy_url = source["proxy_url"];
	        this.use_proxy = source["use_proxy"];
	        this.magic_name = source["magic_name"];
	        this.download_dir = source["download_dir"];
	        this.download_video = source["download_video"];
	        this.download_audio = source["download_audio"];
	        this.download_subtitle = source["download_subtitle"];
	        this.download_combine = source["download_combine"];
	        this.download_limit = source["download_limit"];
	    }
	}

}

export namespace plugin {
	
	export class Plugin {
	    id: string;
	    enable: boolean;
	    settings: {[key: string]: string};
	    manifest_version: number;
	    name: string;
	    description: string;
	    author: string;
	    version: string;
	    homepage: string;
	    color: string;
	    docs_url: string;
	    download_url: string;
	    matches: string[];
	    type: string;
	    location: string;
	    state: number;
	    port: number;
	    pid: number;
	
	    static createFrom(source: any = {}) {
	        return new Plugin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.enable = source["enable"];
	        this.settings = source["settings"];
	        this.manifest_version = source["manifest_version"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.version = source["version"];
	        this.homepage = source["homepage"];
	        this.color = source["color"];
	        this.docs_url = source["docs_url"];
	        this.download_url = source["download_url"];
	        this.matches = source["matches"];
	        this.type = source["type"];
	        this.location = source["location"];
	        this.state = source["state"];
	        this.port = source["port"];
	        this.pid = source["pid"];
	    }
	}

}

export namespace proto {
	
	export class Format {
	    id?: string;
	    fid?: number;
	    mime_type?: string;
	    label?: string;
	    code?: string;
	    url?: string;
	
	    static createFrom(source: any = {}) {
	        return new Format(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.fid = source["fid"];
	        this.mime_type = source["mime_type"];
	        this.label = source["label"];
	        this.code = source["code"];
	        this.url = source["url"];
	    }
	}
	export class Segment {
	    mime_type?: string;
	    formats?: Format[];
	
	    static createFrom(source: any = {}) {
	        return new Segment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mime_type = source["mime_type"];
	        this.formats = this.convertValues(source["formats"], Format);
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
	export class Task {
	    id?: string;
	    url?: string;
	    session_id?: string;
	    title?: string;
	    segments?: Segment[];
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.session_id = source["session_id"];
	        this.title = source["title"];
	        this.segments = this.convertValues(source["segments"], Segment);
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
	export class InfoResponse {
	    title?: string;
	    cover?: string;
	    author?: string;
	    tasks?: Task[];
	
	    static createFrom(source: any = {}) {
	        return new InfoResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.cover = source["cover"];
	        this.author = source["author"];
	        this.tasks = this.convertValues(source["tasks"], Task);
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
	export class ParseResponse {
	    id?: string;
	    tasks?: Task[];
	
	    static createFrom(source: any = {}) {
	        return new ParseResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tasks = this.convertValues(source["tasks"], Task);
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

