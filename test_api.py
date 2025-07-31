#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
AHSFNU Media Cloud API 测试脚本
测试所有可用的API接口
"""

import requests
import json
import os
import time
from typing import Dict, Any, Optional

class APITester:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.token = None
        self.user_id = None
        
    def print_response(self, response: requests.Response, title: str):
        """打印响应结果"""
        print(f"\n{'='*50}")
        print(f"测试: {title}")
        print(f"状态码: {response.status_code}")
        print(f"响应头: {dict(response.headers)}")
        try:
            print(f"响应体: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
        except:
            print(f"响应体: {response.text}")
        print(f"{'='*50}\n")
    
    def test_health_check(self):
        """测试服务器健康状态"""
        try:
            response = self.session.get(f"{self.base_url}/")
            self.print_response(response, "服务器健康检查")
            return response.status_code == 200
        except requests.exceptions.ConnectionError:
            print("❌ 无法连接到服务器，请确保服务器正在运行")
            return False
    
    def test_register(self, username: str, email: str, password: str, invite_code: str):
        """测试用户注册"""
        url = f"{self.base_url}/api/v1/auth/register"
        data = {
            "username": username,
            "email": email,
            "password": password,
            "invite_code": invite_code
        }
        
        response = self.session.post(url, json=data)
        self.print_response(response, f"用户注册 - {username}")
        
        if response.status_code == 201:
            result = response.json()
            self.token = result.get('token')
            self.user_id = result.get('user', {}).get('id')
            print(f"✅ 注册成功，用户ID: {self.user_id}")
            return True
        else:
            print(f"❌ 注册失败: {response.json().get('error', '未知错误')}")
            return False
    
    def test_login(self, username: str, password: str):
        """测试用户登录"""
        url = f"{self.base_url}/api/v1/auth/login"
        data = {
            "username": username,
            "password": password
        }
        
        response = self.session.post(url, json=data)
        self.print_response(response, f"用户登录 - {username}")
        
        if response.status_code == 200:
            result = response.json()
            self.token = result.get('token')
            self.user_id = result.get('user', {}).get('id')
            print(f"✅ 登录成功，用户ID: {self.user_id}")
            return True
        else:
            print(f"❌ 登录失败: {response.json().get('error', '未知错误')}")
            return False
    
    def test_upload_material(self, file_path: str, workflow_id: Optional[int] = None):
        """测试上传素材"""
        if not self.token:
            print("❌ 请先登录获取token")
            return False
        
        url = f"{self.base_url}/api/v1/materials"
        headers = {"Authorization": f"Bearer {self.token}"}
        
        # 准备文件数据
        with open(file_path, 'rb') as f:
            files = {'file': f}
            data = {}
            if workflow_id:
                data['workflow_id'] = str(workflow_id)
            
            response = self.session.post(url, headers=headers, files=files, data=data)
        
        self.print_response(response, f"上传素材 - {os.path.basename(file_path)}")
        
        if response.status_code == 200:
            result = response.json()
            material_id = result.get('data', {}).get('id')
            print(f"✅ 上传成功，素材ID: {material_id}")
            return material_id
        else:
            print(f"❌ 上传失败: {response.json().get('error', '未知错误')}")
            return None
    
    def test_update_material(self, material_id: int, update_data: Dict[str, Any]):
        """测试更新素材"""
        if not self.token:
            print("❌ 请先登录获取token")
            return False
        
        url = f"{self.base_url}/api/v1/materials/{material_id}"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }
        
        response = self.session.put(url, headers=headers, json=update_data)
        self.print_response(response, f"更新素材 - ID: {material_id}")
        
        if response.status_code == 200:
            print(f"✅ 更新成功")
            return True
        else:
            print(f"❌ 更新失败: {response.json().get('error', '未知错误')}")
            return False
    
    def test_unauthorized_access(self):
        """测试未授权访问"""
        url = f"{self.base_url}/api/v1/materials"
        response = self.session.get(url)
        self.print_response(response, "未授权访问测试")
        
        if response.status_code == 401:
            print("✅ 未授权访问被正确拒绝")
            return True
        else:
            print("❌ 未授权访问未被正确拒绝")
            return False
    
    def test_invalid_token(self):
        """测试无效token"""
        url = f"{self.base_url}/api/v1/materials"
        headers = {"Authorization": "Bearer invalid_token"}
        response = self.session.get(url, headers=headers)
        self.print_response(response, "无效token测试")
        
        if response.status_code == 401:
            print("✅ 无效token被正确拒绝")
            return True
        else:
            print("❌ 无效token未被正确拒绝")
            return False

def create_test_file(filename: str = "test.txt", content: str = "这是一个测试文件"):
    """创建测试文件"""
    with open(filename, 'w', encoding='utf-8') as f:
        f.write(content)
    return filename

def main():
    """主测试函数"""
    print("🚀 AHSFNU Media Cloud API 测试开始")
    print("="*60)
    
    # 初始化测试器
    tester = APITester()
    
    # 测试服务器连接
    if not tester.test_health_check():
        return
    
    # 测试未授权访问
    tester.test_unauthorized_access()
    tester.test_invalid_token()
    
    # 测试用户注册
    print("\n📝 测试用户注册...")
    test_username = f"testuser_{int(time.time())}"
    test_email = f"{test_username}@example.com"
    test_password = "123456"
    test_invite_code = "TEST123"  # 需要有效的邀请码
    
    register_success = tester.test_register(test_username, test_email, test_password, test_invite_code)
    
    if not register_success:
        print("⚠️  注册失败，尝试登录现有用户...")
        # 尝试登录
        login_success = tester.test_login("admin", "admin123")  # 假设有管理员账户
        if not login_success:
            print("❌ 无法登录，跳过后续测试")
            return
    
    # 创建测试文件
    test_file = create_test_file()
    
    try:
        # 测试上传素材
        print("\n📤 测试上传素材...")
        material_id = tester.test_upload_material(test_file, workflow_id=1)
        
        if material_id:
            # 测试更新素材
            print("\n✏️  测试更新素材...")
            update_data = {
                "original_filename": "更新后的文件名.txt",
                "is_starred": True,
                "is_public": False
            }
            tester.test_update_material(material_id, update_data)
        
        # 测试上传图片
        print("\n🖼️  测试上传图片...")
        # 创建一个简单的测试图片（这里用文本文件模拟）
        test_image = create_test_file("test_image.txt", "模拟图片内容")
        tester.test_upload_material(test_image)
        
    finally:
        # 清理测试文件
        for file in [test_file, "test_image.txt"]:
            if os.path.exists(file):
                os.remove(file)
    
    print("\n🎉 API测试完成！")
    print("="*60)

if __name__ == "__main__":
    main() 