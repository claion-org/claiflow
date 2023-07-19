# webhook payload

```golang
map[string]interface{}{
    "template_uuid": string,
    "cluster_uuid": string,
    "service_uuid": string,
    "service_name": string,
    "inputs": map[string]interface{},
    "assigned_client_uuid": string,
    "status": int,
    "status_description": string,
    "step_count": int,
    "step_position": int,
    "step_status": int,
    "step_started": time.Time,
    "step_ended": time.Time,
    "result": json.Rawmessage,
    "error": string,
    "webhook": string,
}
```
