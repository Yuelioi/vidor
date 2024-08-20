export namespace app {
	
	export class PluginConfig {
	    manifest_version: number;
	    name: string;
	    description: string;
	    author: string;
	    version: string;
	    url: string;
	    docs_url: string;
	    download_url: string;
	    matches: string[];
	    settings: string[];
	
	    static createFrom(source: any = {}) {
	        return new PluginConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.manifest_version = source["manifest_version"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.version = source["version"];
	        this.url = source["url"];
	        this.docs_url = source["docs_url"];
	        this.download_url = source["download_url"];
	        this.matches = source["matches"];
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
	export class Config {
	    system: SystemConfig;
	    plugins: PluginConfig[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.system = this.convertValues(source["system"], SystemConfig);
	        this.plugins = this.convertValues(source["plugins"], PluginConfig);
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

export namespace proto {
	
	export class Format {
	    id?: number;
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
	        this.mime_type = source["mime_type"];
	        this.label = source["label"];
	        this.code = source["code"];
	        this.url = source["url"];
	    }
	}
	export class Stream {
	    mime_type?: string;
	    formats?: Format[];
	
	    static createFrom(source: any = {}) {
	        return new Stream(source);
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
	export class StreamInfo {
	    id?: string;
	    url?: string;
	    session_id?: string;
	    title?: string;
	    streams?: Stream[];
	
	    static createFrom(source: any = {}) {
	        return new StreamInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.session_id = source["session_id"];
	        this.title = source["title"];
	        this.streams = this.convertValues(source["streams"], Stream);
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
	    stream_infos?: StreamInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ParseResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.stream_infos = this.convertValues(source["stream_infos"], StreamInfo);
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
	export class ShowResponse {
	    title?: string;
	    cover?: string;
	    author?: string;
	    stream_infos?: StreamInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ShowResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.cover = source["cover"];
	        this.author = source["author"];
	        this.stream_infos = this.convertValues(source["stream_infos"], StreamInfo);
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

export namespace shared {
	
	export class Part {
	
	
	    static createFrom(source: any = {}) {
	        return new Part(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

