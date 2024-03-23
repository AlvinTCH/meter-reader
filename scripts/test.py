import requests
import os


class Test:
  def __init__(self):
    self.run()
    self.test_bad_file()
    self.test_update_file()

  def run(self):
    with open('sample_input.csv', 'rb') as f:
      response = requests.post(
        'http://localhost:8080/upload-csv',
        files={
          'file': f
        }, 
      )
      if response.status_code != 200:
        print('Upload file test failed')
        return
      
      print('Upload file test passed')

  def test_bad_file(self):
    with open('bad_file.csv', 'rb') as f:
      response = requests.post(
        'http://localhost:8080/upload-csv',
        files={
          'file': f
        }, 
      )
      if response.status_code != 400:
        print('Upload bad file test failed')
        return
      print('Upload bad file test passed')

  def test_update_file(self):
    with open('test_update.csv', 'rb') as f:
      response = requests.post(
        'http://localhost:8080/upload-csv',
        files={
          'file': f
        }, 
      )
      if response.status_code != 200:
        print('Upload file test failed')
        return
    
    get_data_res = requests.get('http://localhost:8080/get-update-readings')
    if get_data_res.status_code != 200:
        print('Update test failed')
        return
    
    data = get_data_res.json()
    if len(data) != 1:
        print('Update test failed')
        return
    
    if data[0]['reading'] != 1.69:
        print('Update test failed')
        return
    
    print('Update test passed')

Test()