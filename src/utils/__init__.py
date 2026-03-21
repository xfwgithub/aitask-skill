import json
from typing import Any, Optional


def dumps(obj: Optional[Any]) -> Optional[str]:
    """将 Python 对象序列化为 JSON 字符串"""
    if obj is None:
        return None
    return json.dumps(obj, ensure_ascii=False)


def loads(s: Optional[str]) -> Optional[Any]:
    """将 JSON 字符串反序列化为 Python 对象"""
    if s is None:
        return None
    return json.loads(s)
