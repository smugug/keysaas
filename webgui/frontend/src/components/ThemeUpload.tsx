import React, { useState } from 'react';
// import {uploadTheme} from '../api';
import axios from 'axios';
import { THEMES_URL } from '../constants';
import './ThemeUpload.css';

interface ThemeUploadProps {
    name: string;
}

const ThemeUpload: React.FC<ThemeUploadProps> = ({name}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadStatus, setUploadStatus] = useState<string | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      alert('Please select a file first');
      return;
    }
    // const formData = new FormData();
    // formData.append('file', selectedFile);

    const form = new FormData();
    form.append('file', selectedFile);

    //uploadTheme(name,formData);
    await axios({
        method: "post",
        url: `${THEMES_URL}/${name}`,
        data: form,
        headers: { "Content-Type": "multipart/form-data" },
    }).then(function (response) {
      //handle success
      setUploadStatus('File uploaded successfully');
      console.log(response);
    })
    .catch(function (response) {
      //handle error
      setUploadStatus('Error uploading file');
      console.log(response);
    });
  };

  return (
    <div className="themediv">
      <input type="file" onChange={handleFileChange} />
      <button className="themebutton" onClick={handleUpload}>Upload Theme</button>
      {uploadStatus && <p>{uploadStatus}</p>}
    </div>
  );
};

export default ThemeUpload;
