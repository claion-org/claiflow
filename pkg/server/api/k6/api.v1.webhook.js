import http from 'k6/http';

export class Webhook {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    Create(uuid, name, summary, url, method, headers, timeout, conditionValidator, conditionFilter) {
        const path = `${this.baseURL}/api/v1/webhook`;
        const payload = JSON.stringify({
            uuid: uuid,
            name: name,
            summary: summary,
            url: url, 
            method: method, 
            headers: headers, 
            timeout: timeout, 
            conditionValidator: conditionValidator,
            conditionFilter: conditionFilter,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.post(path, payload, params);
    }

    Get(uuid) {
        const path = `${this.baseURL}/api/v1/webhook/${uuid}`;
    
        return http.get(path);
    }

    Update(uuid, name, summary, url, method, headers, timeout, conditionValidator, conditionFilter) {
        const path = `${this.baseURL}/api/v1/webhook/${uuid}`;
        const payload = JSON.stringify({
            name: name,
            summary: summary,
            url: url, 
            method: method, 
            headers: headers, 
            timeout: timeout, 
            conditionValidator: conditionValidator,
            conditionFilter: conditionFilter,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.put(path, payload, params);
    }

    Delete(uuid) {
        const path = `${this.baseURL}/api/v1/webhook/${uuid}`;
    
        return http.del(path)
    }
}
