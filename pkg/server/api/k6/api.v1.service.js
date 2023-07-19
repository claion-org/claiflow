import http from 'k6/http';

export class Service {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    Create(cluster_uuid, uuid, name, summary, template_uuid, inputs, subscribed_channel) {
        const path = `${this.baseURL}/api/v1/cluster/${cluster_uuid}/service`;
        const payload = JSON.stringify({
            uuid: uuid,
            name: name,
            summary: summary,
            template_uuid: template_uuid,
            inputs: inputs,
            template_uuid: template_uuid,
            subscribed_channel: subscribed_channel,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.post(path, payload, params);
    }

    CreateMultiClusters(cluster_uuids, uuid, name, summary, template_uuid, inputs, subscribed_channel) {
        const path = `${this.baseURL}/api/v1/service`;
        const payload = JSON.stringify({
            cluster_uuids: cluster_uuids,
            uuid: uuid,
            name: name,
            summary: summary,
            template_uuid: template_uuid,
            inputs: inputs,
            template_uuid: template_uuid,
            subscribed_channel: subscribed_channel,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.post(path, payload, params);
    }


    Get(cluster_uuid, uuid) {
        const path = `${this.baseURL}/api/v1/cluster/${cluster_uuid}/service/${uuid}`;
    
        return http.get(path);
    }
}
