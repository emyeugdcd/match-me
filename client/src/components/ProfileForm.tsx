import React, { useState, useEffect } from 'react';
import profileService from '../services/profileService';
import {
  type ProfileUpdateRequest,
  type BioDataUpdateRequest,
  type ProfileFormState,
  type BioFormState,
  LOCATIONS,
} from '../types';

interface ProfileFormProps {
  initialProfile?: ProfileUpdateRequest;
  initialBio?: BioDataUpdateRequest;
  profilePictureUrl?: string;
  onSave?: () => void;
  onCancel?: () => void;
}

type ProfileStringField =
  | 'first_name'
  | 'last_name'
  | 'username'
  | 'about_me'
  | 'location'
  | 'profile_picture_url';

const CITY_COORDINATES: Record<string, { lat: number; lon: number }> = {
  Helsinki: { lat: 60.1699, lon: 24.9384 },
  Espoo: { lat: 60.2055, lon: 24.6559 },
  Vantaa: { lat: 60.2934, lon: 25.0378 },
  Tampere: { lat: 61.4978, lon: 23.7610 },
  Turku: { lat: 60.4518, lon: 22.2666 },
  Oulu: { lat: 65.0121, lon: 25.4651 },
  Jyväskylä: { lat: 62.2426, lon: 25.7473 },
  Lahti: { lat: 60.9827, lon: 25.6612 },
  Kuopio: { lat: 62.8924, lon: 27.6770 },
};

const ProfileForm: React.FC<ProfileFormProps> = ({
  initialProfile,
  initialBio,
  profilePictureUrl,
  onSave,
  onCancel,
}) => {
  const [profileForm, setProfileForm] = useState<ProfileFormState>({
    first_name: '',
    last_name: '',
    username: '',
    about_me: '',
    location: '',
    profile_picture_url: '',
    latitude: null,
    longitude: null,
  });

  const [bioForm, setBioForm] = useState<BioFormState>({
    interests: [''],
    hobbies: [''],
    music_preferences: [''],
    food_preferences: [''],
    looking_for: '',
    age_range_min: undefined,
    age_range_max: undefined,
    max_distance_km: undefined,
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [geoStatus, setGeoStatus] = useState<string | null>(null);

  useEffect(() => {
    if (initialProfile) {
      setProfileForm(prev => ({
        ...prev,
        ...initialProfile,
        latitude:
          initialProfile.latitude !== undefined ? initialProfile.latitude : prev.latitude,
        longitude:
          initialProfile.longitude !== undefined ? initialProfile.longitude : prev.longitude,
      }));
    }

    if (initialBio) {
      setBioForm(prev => ({
        ...prev,
        ...initialBio,
        max_distance_km:
          initialBio.max_distance_km !== undefined
            ? initialBio.max_distance_km
            : prev.max_distance_km,
      }));
    }
  }, [initialProfile, initialBio]);

  useEffect(() => {
    if (profilePictureUrl !== undefined) {
      setProfileForm(prev => ({
        ...prev,
        profile_picture_url: profilePictureUrl,
      }));
    }
  }, [profilePictureUrl]);

  const handleProfileChange = (field: ProfileStringField, value: string) => {
    setProfileForm(prev => {
      const updates: Partial<ProfileFormState> = { [field]: value };

      // Auto-set coordinates if location matches a known city
      if (field === 'location' && CITY_COORDINATES[value]) {
        updates.latitude = CITY_COORDINATES[value].lat;
        updates.longitude = CITY_COORDINATES[value].lon;
      }

      return { ...prev, ...updates };
    });
  };

  const handleBioArrayChange = (field: keyof BioFormState, index: number, value: string) => {
    setBioForm(prev => {
      const newArray = [...(prev[field] as string[])];
      newArray[index] = value;
      return { ...prev, [field]: newArray };
    });
  };

  const addArrayItem = (field: keyof BioFormState) => {
    const currentArray = bioForm[field] as string[];
    // Only add if last item isn't empty
    if (currentArray[currentArray.length - 1]?.trim() !== '') {
      setBioForm(prev => ({
        ...prev,
        [field]: [...currentArray, ''],
      }));
    }
  };

  const removeArrayItem = (field: keyof BioFormState, index: number) => {
    setBioForm(prev => ({
      ...prev,
      [field]: (prev[field] as string[]).filter((_, i) => i !== index),
    }));
  };

  const handleBioChange = <K extends keyof BioFormState>(field: K, value: BioFormState[K]) => {
    setBioForm(prev => ({ ...prev, [field]: value }));
  };

  const handleUseCurrentLocation = () => {
    if (!navigator.geolocation) {
      setError('Your browser does not support location services.');
      return;
    }

    setError(null);
    setGeoStatus('Getting your current location…');

    navigator.geolocation.getCurrentPosition(
      position => {
        const { latitude, longitude } = position.coords;
        setProfileForm(prev => ({
          ...prev,
          latitude,
          longitude,
        }));
        setGeoStatus('Current location saved (approximate).');
      },
      err => {
        console.error(err);
        setError('Could not get your current location.');
        setGeoStatus(null);
      },
      {
        enableHighAccuracy: false,
        timeout: 10000,
      }
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      // Filter out empty strings from arrays
      const cleanedBio: BioDataUpdateRequest = {
        ...bioForm,
        interests: bioForm.interests.filter(item => item.trim() !== ''),
        hobbies: bioForm.hobbies.filter(item => item.trim() !== ''),
        music_preferences: bioForm.music_preferences.filter(item => item.trim() !== ''),
        food_preferences: bioForm.food_preferences.filter(item => item.trim() !== ''),
      };

      // Validate data completeness
      const isComplete = profileService.isProfileDataComplete(profileForm, cleanedBio);

      if (!isComplete) {
        setError('Please fill in all required fields. You need at least 5 bio data points.');
        setLoading(false);
        return;
      }

      if (bioForm.age_range_min && bioForm.age_range_max) {
        if (bioForm.age_range_min > bioForm.age_range_max) {
          setError('Minimum age cannot be greater than maximum age');
          setLoading(false);
          return;
        }
      }

      // Save profile and bio
      await Promise.all([
        profileService.updateMyProfile(profileForm),
        profileService.updateMyBioData(cleanedBio),
      ]);

      setSuccess(true);
      onSave?.();

      // Reset success message after 3 seconds
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save profile');
    } finally {
      setLoading(false);
    }
  };

  const renderArrayField = (
    label: string,
    field: keyof BioFormState,
    placeholder: string
  ) => (
    <div className="bio-array-field">
      <label>{label}</label>
      {(bioForm[field] as string[]).map((item, index) => (
        <div key={index} className="array-input-group">
          <input
            type="text"
            value={item}
            onChange={e => handleBioArrayChange(field, index, e.target.value)}
            placeholder={placeholder}
          />
          {(bioForm[field] as string[]).length > 1 && (
            <button
              type="button"
              onClick={() => removeArrayItem(field, index)}
              className="remove-btn"
            >
              Remove
            </button>
          )}
        </div>
      ))}
      <button type="button" onClick={() => addArrayItem(field)} className="add-btn">
        Add {label.toLowerCase()}
      </button>
    </div>
  );

  return (
    <form onSubmit={handleSubmit} className="profile-form">
      <div className="form-section">
        <h3>Basic Information</h3>

        <div className="form-row">
          <div className="form-group">
            <label>First Name *</label>
            <input
              type="text"
              value={profileForm.first_name}
              onChange={e => handleProfileChange('first_name', e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label>Last Name *</label>
            <input
              type="text"
              value={profileForm.last_name}
              onChange={e => handleProfileChange('last_name', e.target.value)}
              required
            />
          </div>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>Username</label>
            <input
              type="text"
              value={profileForm.username}
              onChange={e => handleProfileChange('username', e.target.value)}
              placeholder="Unique username"
            />
          </div>

          <div className="form-group">
            <label>Location *</label>
            <div className="location-input-group" style={{ display: 'flex', gap: '0.5rem', flexDirection: 'column' }}>
              <select
                value={profileForm.location}
                onChange={e => handleProfileChange('location', e.target.value)}
                style={{ marginBottom: '0.5rem' }}
              >
                <option value="">Select city</option>
                {LOCATIONS.map(city => (
                  <option key={city} value={city}>
                    {city}
                  </option>
                ))}
              </select>
              <input
                type="text"
                value={profileForm.location}
                onChange={e => handleProfileChange('location', e.target.value)}
                placeholder="Or type city..."
                required
              />
            </div>
            <button
              type="button"
              onClick={handleUseCurrentLocation}
              className="secondary-btn"
              style={{ marginTop: '0.5rem', fontSize: '0.85rem', padding: '0.25rem 0.5rem', background: 'transparent', border: '1px solid var(--border-color)', color: 'var(--text-secondary)', borderRadius: '4px', cursor: 'pointer' }}
            >
              📍 Use current location
            </button>
            {profileForm.latitude !== null && profileForm.longitude !== null && (
              <p className="field-hint">
                Coords: {profileForm.latitude.toFixed(2)}, {profileForm.longitude.toFixed(2)}
              </p>
            )}
            {geoStatus && <p className="field-hint">{geoStatus}</p>}
          </div>
        </div>

        <div className="form-group">
          <label>About Me</label>
          <textarea
            value={profileForm.about_me}
            onChange={e => handleProfileChange('about_me', e.target.value)}
            rows={4}
            placeholder="Tell others about yourself..."
          />
        </div>
      </div>

      <div className="form-section">
        <h3>Bio Information</h3>

        {renderArrayField('Interests', 'interests', 'e.g., Reading, Hiking, Gaming')}
        {renderArrayField('Programming Languages', 'hobbies', 'e.g., Python, Go, TypeScript')}
        {renderArrayField(
          'Music Preferences',
          'music_preferences',
          'e.g., Rock, Jazz, Classical'
        )}
        {renderArrayField(
          'Food Preferences',
          'food_preferences',
          'e.g., Italian, Vegan, BBQ'
        )}

        <div className="form-group">
          <label>Looking For</label>
          <textarea
            value={bioForm.looking_for}
            onChange={e => handleBioChange('looking_for', e.target.value)}
            placeholder="What are you looking for in connections?"
            rows={3}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>Age Range Min</label>
            <input
              type="number"
              value={bioForm.age_range_min ?? ''}
              onChange={e =>
                handleBioChange(
                  'age_range_min',
                  e.target.value ? parseInt(e.target.value, 10) : undefined
                )
              }
              min="18"
              max="100"
            />
          </div>

          <div className="form-group">
            <label>Age Range Max</label>
            <input
              type="number"
              value={bioForm.age_range_max ?? ''}
              onChange={e =>
                handleBioChange(
                  'age_range_max',
                  e.target.value ? parseInt(e.target.value, 10) : undefined
                )
              }
              min="18"
              max="100"
            />
          </div>
        </div>

        <div className="form-group">
          <label>Maximum distance (km) for matches</label>
          <input
            type="number"
            min={1}
            max={500}
            value={bioForm.max_distance_km ?? ''}
            onChange={e =>
              handleBioChange(
                'max_distance_km',
                e.target.value ? parseInt(e.target.value, 10) : undefined
              )
            }
            placeholder="e.g., 10, 25, 50"
          />
          <p className="field-hint">
            We’ll prefer people within this distance from your location.
          </p>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">Profile saved successfully!</div>}

      <div className="form-actions">
        <button
          type="button"
          onClick={onCancel}
          disabled={loading}
          className="cancel-btn"
        >
          Cancel
        </button>
        <button type="submit" disabled={loading} className="submit-btn">
          {loading ? 'Saving...' : 'Save Profile'}
        </button>
      </div>
    </form>
  );
};

export default ProfileForm;
