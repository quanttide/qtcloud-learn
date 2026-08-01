use clap::{Parser, Subcommand};

use qtcloud_learn_cli::{api, commands};

/// 量潮学习云 CLI（学员侧）
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
    /// 学员管理
    Student {
        #[command(subcommand)]
        cmd: commands::student::StudentCmd,
    },
    /// 班级管理
    Class {
        #[command(subcommand)]
        cmd: commands::class::ClassCmd,
    },
    /// 选课 / 报名
    Enrollment {
        #[command(subcommand)]
        cmd: commands::enrollment::EnrollmentCmd,
    },
    /// 学习进度
    Progress {
        #[command(subcommand)]
        cmd: commands::progress::ProgressCmd,
    },
    /// 考核管理
    Assessment {
        #[command(subcommand)]
        cmd: commands::assessment::AssessmentCmd,
    },
}

fn main() {
    let cli = Cli::parse();
    let api = api::ApiClient::new(&cli.base_url);

    let output = match cli.command {
        Commands::Version => Ok(format!(
            "qtcloud-learn {}",
            env!("CARGO_PKG_VERSION")
        )),
        Commands::Student { cmd } => commands::student::run(&api, cmd),
        Commands::Class { cmd } => commands::class::run(&api, cmd),
        Commands::Enrollment { cmd } => commands::enrollment::run(&api, cmd),
        Commands::Progress { cmd } => commands::progress::run(&api, cmd),
        Commands::Assessment { cmd } => commands::assessment::run(&api, cmd),
    };

    match output {
        Ok(text) => println!("{text}"),
        Err(e) => {
            eprintln!("错误：{e}");
            std::process::exit(1);
        }
    }
}
