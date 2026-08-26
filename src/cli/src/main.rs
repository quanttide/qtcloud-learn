use clap::{Parser, Subcommand};

use qtcloud_learn_cli::{api, commands};

/// 量潮学习云 CLI（对齐《量潮学习管理标准》：Learner × Lesson → Completion）
#[derive(Parser)]
#[command(name = "qtcloud-learn", version, about)]
struct Cli {
    /// Provider API 地址
    #[arg(long, global = true, default_value = "http://localhost:8080")]
    base_url: String,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// 打印版本信息
    Version,
    /// 学习者管理
    Learner {
        #[command(subcommand)]
        cmd: commands::learner::LearnerCmd,
    },
    /// 完成记录管理
    Completion {
        #[command(subcommand)]
        cmd: commands::completion::CompletionCmd,
    },
}

fn main() {
    let cli = Cli::parse();
    let api = api::ApiClient::new(&cli.base_url);

    let output = match cli.command {
        Commands::Version => Ok(format!("qtcloud-learn {}", env!("CARGO_PKG_VERSION"))),
        Commands::Learner { cmd } => commands::learner::run(&api, cmd),
        Commands::Completion { cmd } => commands::completion::run(&api, cmd),
    };

    match output {
        Ok(text) => println!("{text}"),
        Err(e) => {
            eprintln!("错误：{e}");
            std::process::exit(1);
        }
    }
}
