import requests
import json
import os


def fetch_and_save_models_map():
    url = "https://openrouter.ai/api/v1/models"

    try:
        response = requests.get(url, timeout=15)
        response.raise_for_status()
        raw_data = response.json().get("data", [])
    except Exception as e:
        print(f"出错: {e}")
        return

    models_map = {}

    for item in raw_data:
        model_id = item.get("id")
        if not model_id: continue
        provider = model_id.split("/")[0]
        model_id = model_id.split(":")[0].replace(provider+"/","")

        model_name = item.get("name", model_id)
        description = item.get("description") or ""
        arch = item.get("architecture") or {}
        input_modalities = arch.get("input_modalities", [])
        output_modalities = arch.get("output_modalities", [])

        context_length = item.get("context_length", 0)

        models_map[model_id] = {
            "provider": provider,
            "name": model_name,
            "description": description,
            "supports_image": "image" in input_modalities,
            "context_length": context_length,
            "inputs": input_modalities,
            "outputs": output_modalities,
        }

    # 输出到 frontend/src/config/model-capabilities.json
    script_dir = os.path.dirname(os.path.abspath(__file__))
    repo_root = os.path.dirname(script_dir)
    output_file = os.path.join(repo_root, "frontend", "src", "config", "model-capabilities.json")
    
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(models_map, f, ensure_ascii=False, indent=2)
    
    print(f"✅ 已更新 {len(models_map)} 个模型到 {output_file}")


if __name__ == "__main__":
    fetch_and_save_models_map()
