import { AtSign, Clock } from "lucide-react";

const mentions = [
  {
    id: 1,
    author: "@techwriter",
    authorAvatar: 1,
    text: "Great insights from @yourhandle on building modern web applications. Definitely worth a read!",
    time: "2h ago",
  },
  {
    id: 2,
    author: "@designpro",
    authorAvatar: 2,
    text: "Shoutout to @yourhandle for the amazing glassmorphism tutorial. My designs have leveled up!",
    time: "4h ago",
  },
  {
    id: 3,
    author: "@devguru",
    authorAvatar: 3,
    text: "Following @yourhandle's advice on responsive design testing. Already caught 3 bugs!",
    time: "6h ago",
  },
  {
    id: 4,
    author: "@codemaster",
    authorAvatar: 4,
    text: "@yourhandle Do you have any recommendations for state management in large React apps?",
    time: "7h ago",
  },
  {
    id: 5,
    author: "@webdev101",
    authorAvatar: 5,
    text: "Just saw @yourhandle's new feature launch. The attention to detail is incredible!",
    time: "9h ago",
  },
  {
    id: 6,
    author: "@uiuxdesigner",
    authorAvatar: 6,
    text: "@yourhandle Your dark mode take resonates so much. Light mode hurts my eyes 😅",
    time: "10h ago",
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
  ];
  return colors[id % colors.length];
};

export default function MentionsToday() {
  return (
    <div className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-6 h-fit">
      <div className="flex items-center gap-2 mb-4">
        <AtSign className="w-5 h-5 text-orange-400" />
        <h2 className="text-xl font-bold text-white">Mentions Today</h2>
      </div>
      <div className="space-y-4">
        {mentions.map((mention) => (
          <div key={mention.id} className="pb-4 border-b border-white/10 last:border-0 last:pb-0">
            <div className="flex items-start gap-3">
              <div
                className={`w-10 h-10 rounded-full bg-gradient-to-br ${getAvatarColor(
                  mention.authorAvatar
                )} flex-shrink-0`}
              />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-white font-medium">{mention.author}</span>
                  <span className="text-xs text-gray-500">·</span>
                  <div className="flex items-center gap-1 text-gray-400 text-xs">
                    <Clock className="w-3 h-3" />
                    <span>{mention.time}</span>
                  </div>
                </div>
                <p className="text-sm text-gray-300 leading-relaxed">{mention.text}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
