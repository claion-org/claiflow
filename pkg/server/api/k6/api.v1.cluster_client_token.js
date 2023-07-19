import http from 'k6/http';

export class ClusterClientToken {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    Create(clusterUuid, uuid, name, summary, token) {
        const path = `${this.baseURL}/api/v1/cluster_token`;
        const payload = JSON.stringify({
            cluster_uuid: clusterUuid,
            uuid: uuid,
            name: name,
            summary: summary,
            token: token,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.post(path, payload, params);
    }

    Get(uuid) {
        const path = `${this.baseURL}/api/v1/cluster_token/${uuid}`;
    
        return http.get(path);
    }

    Update(uuid, name, summary, token, issed, exp) {
        const path = `${this.baseURL}/api/v1/cluster_token/${uuid}`;
        const payload = JSON.stringify({
            name: name,
            summary: summary,
            token: token,
            issued_at_time: issed,
            expiration_time: exp,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.put(path, payload, params);
    }

    Refresh(uuid) {
        const path = `${this.baseURL}/api/v1/cluster_token/${uuid}/refresh`;
    
        return http.put(path);
    }

    Expire(uuid) {
        const path = `${this.baseURL}/api/v1/cluster_token/${uuid}/expire`;
    
        return http.put(path);
    }

    Delete(uuid) {
        const path = `${this.baseURL}/api/v1/cluster_token/${uuid}`;
    
        return http.del(path)
    }
}
