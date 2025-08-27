export type UserProfile = {
   communication_preference: string;
  id: number;
  user_id: number;
  height: number;
  weight: number;
  age: number;
  gender: string;
  goals: string;
  allergies: string;
  medical_history: string;
  medications: string;
  physio_needs: string;
  preferred_workout_time?: string;
  profile_image_url?: string;
  target_weight: number;
  timeline: number;
  workout_days: string;
  is_profile_complete: boolean;
};

export type UserProfileResponse = {
  profile: UserProfile;
  is_complete: boolean;
  completion_percentage: number;
};
