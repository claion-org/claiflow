import http from 'k6/http';

export class UUIDs {
    ClusterTokenExpirationTime() { return '0f5658f37f2b45d881f19c7f56ea2e23'}
    ClientSessionSignatureSecret() { return '77f7b2aeb0aa4254ad073ae7743291ab' }
    ClientSessionExpirationTime() { return  'af9a14a58b254d13ae69c065a27811b6' }
    ClientConfigServiceValidityPeriod() { return 'bc2cd0f95b6d4db68870d30862523a04' }	
}


export class GlobalVariables {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    Get(uuid) {
        const path = `${this.baseURL}/api/v1/global_variables/${uuid}`;
    
        return http.get(path);
    }

    Find() {
        const path = `${this.baseURL}/api/v1/global_variables`;
    
        return http.get(path);
    }

    Update(uuid, value) {
        const path = `${this.baseURL}/api/v1/global_variables/${uuid}`;
        const payload = JSON.stringify({
            value: value,
        });
        const params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        return http.put(path, payload, params);
    }
}
