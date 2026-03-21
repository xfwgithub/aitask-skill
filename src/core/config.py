from pydantic_settings import BaseSettings
from typing import Optional
import os


class Settings(BaseSettings):
    """应用配置"""
    
    # 基础配置
    APP_NAME: str = "任务管理系统"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = True
    
    # 数据库配置
    DATABASE_URL: str = "sqlite+aiosqlite:///./data/tasks.db"
    
    # API 配置
    API_PREFIX: str = "/api/v1"
    
    # 服务器配置
    SERVER_PORT: int = 8000
    SERVER_HOST: str = "0.0.0.0"
    
    # 安全配置
    SECRET_KEY: str = "your-secret-key-change-in-production"
    ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 7  # 7 天
    
    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()
