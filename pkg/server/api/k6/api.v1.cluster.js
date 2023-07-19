import http from 'k6/http';

export class Cluster {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    Create(uuid, name, summary) {
        const path = `${this.baseURL}/api/v1/cluster`;
        const payload = JSON.stringify({
            uuid: uuid,
            name: name,
            summary: summary,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.post(path, payload, params);
    }

    Get(uuid) {
        const path = `${this.baseURL}/api/v1/cluster/${uuid}`;
    
        return http.get(path);
    }

    Update(uuid, name, summary) {
        const path = `${this.baseURL}/api/v1/cluster/${uuid}`;
        const payload = JSON.stringify({
            name: name,
            summary: summary,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.put(path, payload, params);
    }

    Delete(uuid) {
        const path = `${this.baseURL}/api/v1/cluster/${uuid}`;
    
        return http.del(path)
    }
}
