#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速API测试脚本
"""

import requests
import json
import os

def test_api():
    base_url = "http://localhost:8080"
    
    print("🚀 快速API测试")
    print("="*40)
    
    # 1. 测试服务器连接
    try:
        response = requests.get(f"{base_url}/")
        print(f"✅ 服务器连接正常 - 状态码: {response.status_code}")
    except:
        print("❌ 无法连接到服务器")
        return
    
    # 2. 测试注册
    print("\n📝 测试用户注册...")
    register_data = {
        "username": "testuser123",
        "email": "test@example.com", 
        "password": "123456",
        "invite_code": "TEST123"  # 需要有效的邀请码
    }
    
    response = requests.post(f"{base_url}/api/v1/auth/register", json=register_data)
    print(f"注册响应: {response.status_code}")
    if response.status_code == 201:
        result = response.json()
        token = result.get('token')
        print(f"✅ 注册成功，获得token: {token[:20]}...")
    else:
        print(f"❌ 注册失败: {response.text}")
        # 尝试登录
        print("\n🔐 尝试登录...")
        login_data = {
            "username": "testuser123",
            "password": "123456"
        }
        response = requests.post(f"{base_url}/api/v1/auth/login", json=login_data)
        if response.status_code == 200:
            result = response.json()
            token = result.get('token')
            print(f"✅ 登录成功，获得token: {token[:20]}...")
        else:
            print(f"❌ 登录失败: {response.text}")
            return
    
    # 3. 测试上传文件
    print("\n📤 测试文件上传...")
    
    # 创建测试文件
    test_content = "这是一个测试文件内容"
    with open("test_upload.txt", "w", encoding="utf-8") as f:
        f.write(test_content)
    
    try:
        headers = {"Authorization": f"Bearer {token}"}
        with open(r"C:\Users\fangk\Desktop\14班照片\曹景添.jpg", "rb") as f:
            files = {"file": f}
            data = {"workflow_id": "1"}
            
            response = requests.post(
                f"{base_url}/api/v1/materials",
                headers=headers,
                files=files,
                data=data
            )
        
        print(f"上传响应: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            material_id = result.get('data', {}).get('id')
            print(f"✅ 上传成功，素材ID: {material_id}")
            5
            # 4. 测试更新素材
            print("\n✏️ 测试更新素材...")
            update_data = {
                "original_filename": "更新后的文件名.txt",
                "is_starred": True,
                "is_public": False
            }
            
            response = requests.put(
                f"{base_url}/api/v1/materials/{material_id}",
                headers=headers,
                json=update_data
            )
            
            print(f"更新响应: {response.status_code}")
            if response.status_code == 200:
                print("✅ 更新成功")
            else:
                print(f"❌ 更新失败: {response.text}")
        else:
            print(f"❌ 上传失败: {response.text}")
    
    finally:
        # 清理测试文件
        if os.path.exists("test_upload.txt"):
            os.remove("test_upload.txt")
    
    print("\n🎉 测试完成！")

if __name__ == "__main__":
    test_api() 