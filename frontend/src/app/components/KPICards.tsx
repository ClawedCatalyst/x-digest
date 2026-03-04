import { FileText, Heart, MessageCircle, AtSign, TrendingUp } from "lucide-react";

const kpiData = [
  {
    label: "Posts Today",
    value: "5",
    change: "+2",
    icon: FileText,
    gradient: "from-blue-500 to-cyan-500",
  },
  {
    label: "New Likes",
    value: "127",
    change: "+34",
    icon: Heart,
    gradient: "from-pink-500 to-rose-500",
  },
  {
    label: "New Replies",
    value: "23",
    change: "+8",
    icon: MessageCircle,
    gradient: "from-purple-500 to-indigo-500",
  },
  {
    label: "Mentions",
    value: "12",
    change: "+5",
    icon: AtSign,
    gradient: "from-orange-500 to-yellow-500",
  },
  {
    label: "Top Engager Score",
    value: "94",
    change: "+12",
    icon: TrendingUp,
    gradient: "from-green-500 to-emerald-500",
  },
];

export default function KPICards() {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
      {kpiData.map((kpi, index) => {
        const Icon = kpi.icon;
        return (
          <div
            key={index}
            className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-4 lg:p-5 hover:bg-white/10 transition-all"
          >
            <div className="flex items-start justify-between mb-3">
              <div className={`w-10 h-10 rounded-[12px] bg-gradient-to-br ${kpi.gradient} flex items-center justify-center`}>
                <Icon className="w-5 h-5 text-white" />
              </div>
              <span className="text-xs text-green-400 bg-green-400/10 px-2 py-1 rounded-full">
                {kpi.change}
              </span>
            </div>
            <div>
              <p className="text-2xl lg:text-3xl font-bold text-white mb-1">{kpi.value}</p>
              <p className="text-xs lg:text-sm text-gray-400">{kpi.label}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
