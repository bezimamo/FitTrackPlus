type Props = {
  name: string;
  profileImage?: string;
  isComplete: boolean;
};

export default function ProfileHeader({ name, profileImage, isComplete }: Props) {
  return (
    <div className="relative w-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-3xl shadow-2xl p-8 flex items-center gap-8 text-white overflow-hidden">
      
      {/* Decorative background shapes */}
      <div className="absolute top-0 left-0 w-40 h-40 bg-white/10 rounded-full blur-3xl -z-10"></div>
      <div className="absolute bottom-0 right-0 w-56 h-56 bg-white/10 rounded-full blur-4xl -z-10"></div>

      {/* Profile Image */}
      {profileImage ? (
        <img
          src={profileImage}
          alt={name}
          className="w-24 h-24 rounded-full border-4 border-white shadow-xl object-cover"
        />
      ) : (
        <div className="w-24 h-24 rounded-full border-4 border-white shadow-xl flex items-center justify-center bg-white/20 text-4xl font-bold">
          {name.charAt(0)}
        </div>
      )}

      {/* Profile Info */}
      <div className="flex flex-col justify-center">
        <h2 className="text-3xl font-extrabold drop-shadow-lg">{name}</h2>
        <p
          className={`mt-2 text-lg font-semibold ${
            isComplete ? "text-green-200" : "text-red-200"
          }`}
        >
          {isComplete ? "✔ Profile Complete" : "✘ Profile Incomplete"}
        </p>

        {/* Extra Space for future info */}
        <div className="mt-2 text-sm text-white/80">
          Keep your profile updated for better recommendations!
        </div>
      </div>
    </div>
  );
}
