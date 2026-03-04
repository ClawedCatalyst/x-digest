import { Trophy, Heart, MessageCircle, AtSign, TrendingUp } from "lucide-react";

const topEngagers = [
  {
    id: 1,
    rank: 1,
    username: "@sarahsmith",
    avatar: 1,
    score: 94,
    engagement: {
      likes: 23,
      replies: 8,
      mentions: 2,
    },
  },
  {
    id: 2,
    rank: 2,
    username: "@johndoe",
    avatar: 2,
    score: 87,
    engagement: {
      likes: 19,
      replies: 6,
      mentions: 3,
    },
  },
  {
    id: 3,
    rank: 3,
    username: "@mikejones",
    avatar: 3,
    score: 76,
    engagement: {
      likes: 15,
      replies: 5,
      mentions: 1,
    },
  },
  {
    id: 4,
    rank: 4,
    username: "@emilybrown",
    avatar: 4,
    score: 68,
    engagement: {
      likes: 12,
      replies: 4,
      mentions: 2,
    },
  },
  {
    id: 5,
    rank: 5,
    username: "@alexwilson",
    avatar: 5,
    score: 61,
    engagement: {
      likes: 11,
      replies: 3,
      mentions: 1,
    },
  },
];

const getAvatarColor = (id: number) => {
  const colors = [
    "from-blue-500 to-cyan-500",
    "from-purple-500 to-pink-500",
    "from-orange-500 to-red-500",
    "from-green-500 to-emerald-500",
    "from-indigo-500 to-purple-500",
  ];
  return colors[id % colors.length];
};

const getRankColor = (rank: number) => {
  if (rank === 1) return "from-yellow-500 to-orange-500";
  if (rank === 2) return "from-gray-300 to-gray-400";
  if (rank === 3) return "from-orange-600 to-orange-700";
  return "from-purple-500 to-blue-500";
};

const getRankIcon = (rank: number) => {
  if (rank <= 3) {
    return <Trophy className="w-4 h-4 text-white" />;
  }
  return <span className="text-white font-bold">{rank}</span>;
};

export default function TopEngagersToday() {
  return (
    <div className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-6">
      <div className="flex items-center gap-2 mb-6">
        <TrendingUp className="w-5 h-5 text-green-400" />
        <h2 className="text-xl font-bold text-white">Top Engagers Today</h2>
      </div>
      <div className="space-y-3">
        {topEngagers.map((engager) => (
          <div
            key={engager.id}
            className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-4 hover:bg-white/10 transition-all"
          >
            <div className="flex items-center gap-4">
              {/* Rank badge */}
              <div
                className={`w-10 h-10 rounded-full bg-gradient-to-br ${getRankColor(
                  engager.rank
                )} flex items-center justify-center flex-shrink-0`}
              >
                {getRankIcon(engager.rank)}
              </div>

              {/* Avatar */}
              <div
                className={`w-12 h-12 rounded-full bg-gradient-to-br ${getAvatarColor(
                  engager.avatar
                )} flex-shrink-0`}
              />

              {/* User info and engagement breakdown */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-white font-medium">{engager.username}</span>
                  <div className="flex items-center gap-2 text-green-400">
                    <TrendingUp className="w-4 h-4" />
                    <span className="font-bold">{engager.score}</span>
                  </div>
                </div>
                {/* Engagement stats */}
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-1.5 text-xs">
                    <Heart className="w-3.5 h-3.5 text-pink-400" />
                    <span className="text-gray-300">{engager.engagement.likes}</span>
                  </div>
                  <div className="flex items-center gap-1.5 text-xs">
                    <MessageCircle className="w-3.5 h-3.5 text-blue-400" />
                    <span className="text-gray-300">{engager.engagement.replies}</span>
                  </div>
                  <div className="flex items-center gap-1.5 text-xs">
                    <AtSign className="w-3.5 h-3.5 text-orange-400" />
                    <span className="text-gray-300">{engager.engagement.mentions}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
