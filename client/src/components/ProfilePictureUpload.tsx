import React, { useRef, useState } from 'react';
import profileService from '../services/profileService';

interface ProfilePictureUploadProps {
  currentPictureUrl?: string;
  onUploadSuccess: (newUrl: string) => void;
  onUploadError: (error: string) => void;
  onRemovePicture: () => void;
}

const ProfilePictureUpload: React.FC<ProfilePictureUploadProps> = ({
  currentPictureUrl,
  onUploadSuccess,
  onUploadError,
  onRemovePicture
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (file: File) => {
    // Validate file
    const validationError = profileService.validateImageFile(file);
    if (validationError) {
      onUploadError(validationError);
      return;
    }

    setUploading(true);

    try {
      const result = await profileService.uploadProfilePicture(file);
      onUploadSuccess(result.profile_picture_url);
    } catch (err: any) {
      onUploadError(err.message || 'Failed to upload image');
    } finally {
      setUploading(false);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      handleFileSelect(file);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const file = e.dataTransfer.files?.[0];
    if (file) {
      handleFileSelect(file);
    }
  };

  const handleClickArea = () => {
    fileInputRef.current?.click();
  };

  const handleRemovePicture = async () => {
    try {
      await profileService.deleteProfilePicture();
      onRemovePicture();
    } catch (err: any) {
      onUploadError(err.message || 'Failed to remove profile picture');
    }
  };

  const displayPictureUrl = profileService.getProfilePictureUrl(currentPictureUrl);

  return (
    <div className="profile-picture-upload">
      <div className="current-picture-section">
        {displayPictureUrl ? (
          <div className="picture-preview">
            <img
              src={displayPictureUrl}
              alt="Profile"
              className="profile-picture-preview"
            />
            <button
              onClick={handleRemovePicture}
              disabled={uploading}
              className="remove-picture-btn"
            >
              Remove Picture
            </button>
          </div>
        ) : (
          <div className="no-picture-placeholder">
            <div className="placeholder-icon">👤</div>
            <p>No profile picture</p>
          </div>
        )}
      </div>

      {!displayPictureUrl && (
        <div
          className={`upload-area ${isDragging ? 'dragging' : ''} ${uploading ? 'uploading' : ''}`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={handleClickArea}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
            onChange={handleFileChange}
            style={{ display: 'none' }}
          />

          <div className="upload-content">
            {uploading ? (
              <div className="uploading-indicator">
                <div className="spinner"></div>
                <p>Uploading...</p>
              </div>
            ) : (
              <>
                <div className="upload-icon">📁</div>
                <p className="upload-text">
                  <strong>Click to upload</strong> or drag and drop
                </p>
                <p className="upload-subtext">
                  PNG, JPG, GIF, WebP up to 5MB
                </p>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ProfilePictureUpload;