import { Heart, Users } from "lucide-react";

const likedPosts = [
  {
    id: 1,
    postPreview: "Just launched a new feature for our product...",
    likers: [
      { id: 1, username: "@johndoe", avatar: 1 },
      { id: 2, username: "@sarahsmith", avatar: 2 },
      { id: 3, username: "@mikejones", avatar: 3 },
      { id: 4, username: "@emilybrown", avatar: 4 },
      { id: 5, username: "@alexwilson", avatar: 5 },
    ],
    totalLikes: 32,
    moreCount: 27,
  },
  {
    id: 2,
    postPreview: "Hot take: Dark mode should be the default...",
    likers: [
      { id: 6, username: "@davidlee", avatar: 6 },
      { id: 7, username: "@jessicapark", avatar: 7 },
      { id: 8, username: "@ryanchen", avatar: 8 },
      { id: 9, username: "@lauramartinez", avatar: 9 },
    ],
    totalLikes: 54,
    moreCount: 50,
  },
  {
    id: 3,
    postPreview: "Working on some interesting glassmorphism...",
    likers: [
      { id: 10, username: "@chrisgarcia", avatar: 10 },
      { id: 11, username: "@amandarodriguez", avatar: 11 },
      { id: 12, username: "@kevinwhite", avatar: 12 },
    ],
    totalLikes: 21,
    moreCount: 18,
  },
];

const getAvatarColor = (id: number) => {
  const colors = [
    "from-blue-500 to-cyan-500",
    "from-purple-500 to-pink-500",
    "from-orange-500 to-red-500",
    "from-green-500 to-emerald-500",
    "from-indigo-500 to-purple-500",
    "from-pink-500 to-rose-500",
    "from-yellow-500 to-orange-500",
    "from-teal-500 to-cyan-500",
    "from-red-500 to-pink-500",
    "from-cyan-500 to-blue-500",
    "from-lime-500 to-green-500",
    "from-fuchsia-500 to-purple-500",
  ];
  return colors[id % colors.length];
};

export default function NewLikesToday() {
  return (
    <div className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-6 h-fit">
      <div className="flex items-center gap-2 mb-4">
        <Heart className="w-5 h-5 text-pink-400" />
        <h2 className="text-xl font-bold text-white">New Likes Today</h2>
      </div>
      <div className="space-y-4">
        {likedPosts.map((post) => (
          <div key={post.id} className="pb-4 border-b border-white/10 last:border-0 last:pb-0">
            <p className="text-sm text-gray-400 mb-3 line-clamp-2">{post.postPreview}</p>
            {/* Avatar stack */}
            <div className="flex items-center gap-2 mb-2">
              <div className="flex -space-x-2">
                {post.likers.slice(0, 5).map((liker) => (
                  <div
                    key={liker.id}
                    className={`w-8 h-8 rounded-full bg-gradient-to-br ${getAvatarColor(
                      liker.avatar
                    )} border-2 border-[#141420]`}
                  />
                ))}
                {post.moreCount > 0 && (
                  <div className="w-8 h-8 rounded-full bg-white/10 border-2 border-[#141420] flex items-center justify-center">
                    <span className="text-xs text-gray-300">+{post.moreCount}</span>
                  </div>
                )}
              </div>
            </div>
            {/* Usernames list */}
            <div className="flex flex-wrap gap-1.5 mb-2">
              {post.likers.slice(0, 3).map((liker) => (
                <span key={liker.id} className="text-xs text-gray-300">
                  {liker.username}
                  {liker.id !== post.likers[Math.min(2, post.likers.length - 1)].id && ","}
                </span>
              ))}
              {post.likers.length > 3 && (
                <span className="text-xs text-gray-400">+{post.likers.length - 3} more</span>
              )}
            </div>
            <div className="flex items-center gap-1.5 text-xs text-gray-400">
              <Users className="w-3.5 h-3.5" />
              <span>Up to {post.totalLikes} likers</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
